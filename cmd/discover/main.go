package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Azure/go-autorest/autorest/to"

	"github.com/Azure/azure-sdk-for-go/services/compute/mgmt/2018-06-01/compute"
	"github.com/openshift/openshift-azure/pkg/api"
	"github.com/openshift/openshift-azure/pkg/log"
	"github.com/openshift/openshift-azure/pkg/util/azureclient"
	"github.com/sirupsen/logrus"
)

const fakepubkey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7laRyN4B3YZmVrDEZLZoIuUA72pQ0DpGuZBZWykCofIfCPrFZAJgFvonKGgKJl6FGKIunkZL9Us/mV4ZPkZhBlE7uX83AAf5i9Q8FmKpotzmaxN10/1mcnEE7pFvLoSkwqrQSkrrgSm8zaJ3g91giXSbtqvSIj/vk2f05stYmLfhAwNo3Oh27ugCakCoVeuCrZkvHMaJgcYrIGCuFo6q0Pfk9rsZyriIqEa9AtiUOtViInVYdby7y71wcbl0AbbCZsTSqnSoVxm2tRkOsXV6+8X4SnwcmZbao3H+zfO1GBhQOLxJ4NQbzAa8IJh810rYARNLptgmsd4cYXVOSosTX azureuser"

type Options struct {
	rg       string
	name     string
	location string
	logLevel string
}

var (
	gitCommit = "unknown"
	logger    *logrus.Entry
)

// Resource contains details about an Azure resource.
type Resource struct {
	SubscriptionID string
	ResourceGroup  string
}

// ParseResourceID parses a resource ID into a ResourceDetails struct.
// See https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-template-functions-resource#return-value-4.
func ParseResourceID(resourceID string) (Resource, error) {
	const resourceIDPatternText = `(?i)subscriptions/(.+)/resourceGroups/(.+)`
	resourceIDPattern := regexp.MustCompile(resourceIDPatternText)
	match := resourceIDPattern.FindStringSubmatch(resourceID)
	if len(match) == 0 {
		return Resource{}, fmt.Errorf("parsing failed for %s. Invalid resource Id format", resourceID)
	}
	result := Resource{
		SubscriptionID: match[1],
		ResourceGroup:  match[2],
	}
	return result, nil
}

func GetManagedResourceGroup(ctx context.Context, logger *logrus.Entry, rg, name, l string) (mrg string, err error) {
	authorizer, err := azureclient.NewAuthorizerFromContext(ctx)
	if err != nil {
		return "", err
	}

	appsClient := azureclient.NewApplicationsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})
	logger.Infof("sending request for %s", rg)
	app, err := appsClient.ListByResourceGroup(ctx, rg)
	if err != nil {
		return "", err
	}

	for _, app := range app.Values() {

		if id := *app.ManagedResourceGroupID; id != "" {
			r, err := ParseResourceID(id)
			if err != nil {
				logger.Fatal(err)
			}
			return r.ResourceGroup, nil
		}
	}
	return "", nil
}

func ListVMs(ctx context.Context, logger *logrus.Entry, rg string, vmc azureclient.VirtualMachineScaleSetVMsClient, name string) ([]compute.VirtualMachineScaleSetVM, error) {
	vmPages, err := vmc.List(ctx, rg, name, "", "", "")
	if err != nil {
		return nil, err
	}

	var vms []compute.VirtualMachineScaleSetVM
	for vmPages.NotDone() {
		vms = append(vms, vmPages.Values()...)

		err = vmPages.Next()
		if err != nil {
			return nil, err
		}
	}

	return vms, nil
}

func ListScaleSets(ctx context.Context, logger *logrus.Entry, rg string, ssc azureclient.VirtualMachineScaleSetsClient) ([]compute.VirtualMachineScaleSet, error) {
	vmssPages, err := ssc.List(ctx, rg)
	if err != nil {
		return nil, err
	}

	var scaleSets []compute.VirtualMachineScaleSet
	for vmssPages.NotDone() {
		scaleSets = append(scaleSets, vmssPages.Values()...)

		err = vmssPages.Next()
		if err != nil {
			return nil, err
		}
	}

	return scaleSets, nil
}

func updateInstanceCount(ctx context.Context, logger *logrus.Entry, rg string) []error {
	var errs []error
	authorizer, err := azureclient.NewAuthorizerFromContext(ctx)
	if err != nil {
		return append(errs, err)
	}

	virtualMachineScaleSetVMs := azureclient.NewVirtualMachineScaleSetVMsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})
	virtualMachineScaleSets := azureclient.NewVirtualMachineScaleSetsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})

	logger.Debugf("listing scalesets in rg %s", rg)
	scaleSets, err := ListScaleSets(ctx, logger, rg, virtualMachineScaleSets)
	if err != nil {
		errs = append(errs, err)
	}
	for _, s := range scaleSets {
		logger.Debugf("getting vms in scale set %s", *s.Name)
		vms, err := ListVMs(ctx, logger, rg, virtualMachineScaleSetVMs, *s.Name)
		if err != nil {
			errs = append(errs, err)
		}
		size := len(vms)
		logger.Debugf("resizing %s from %d to %d vms", *s.Name, size, size+1)
		_, err = virtualMachineScaleSets.Update(ctx, rg, *s.Name, compute.VirtualMachineScaleSetUpdate{
			Sku: &compute.Sku{
				Capacity: to.Int64Ptr(int64(size) + 1),
			},
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func updateInstanceType(ctx context.Context, logger *logrus.Entry, rg string) []error {
	var errs []error
	authorizer, err := azureclient.NewAuthorizerFromContext(ctx)
	if err != nil {
		return append(errs, err)
	}

	virtualMachineScaleSets := azureclient.NewVirtualMachineScaleSetsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})

	logger.Debugf("listing scalesets in rg %s", rg)
	scaleSets, err := ListScaleSets(ctx, logger, rg, virtualMachineScaleSets)
	if err != nil {
		errs = append(errs, err)
	}
	for _, s := range scaleSets {
		logger.Debugf("change instance type for %s from %s to %s", *s.Name, api.StandardD4sV3, api.StandardD2sV3)
		_, err = virtualMachineScaleSets.Update(ctx, rg, *s.Name, compute.VirtualMachineScaleSetUpdate{
			Sku: &compute.Sku{
				Name: to.StringPtr(string(api.StandardD2sV3)),
			},
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func updateKey(ctx context.Context, logger *logrus.Entry, rg string) []error {
	var errs []error
	authorizer, err := azureclient.NewAuthorizerFromContext(ctx)
	if err != nil {
		return append(errs, err)
	}

	virtualMachineScaleSets := azureclient.NewVirtualMachineScaleSetsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})

	var sshKeyData = fakepubkey

	logger.Debugf("listing scalesets in rg %s", rg)
	scaleSets, err := ListScaleSets(ctx, logger, rg, virtualMachineScaleSets)
	if err != nil {
		errs = append(errs, err)
	}
	for _, s := range scaleSets {
		logger.Debugf("setting ssh key on scale set %s", *s.Name)
		_, err := virtualMachineScaleSets.Update(ctx, rg, *s.Name, compute.VirtualMachineScaleSetUpdate{
			VirtualMachineScaleSetUpdateProperties: &compute.VirtualMachineScaleSetUpdateProperties{
				VirtualMachineProfile: &compute.VirtualMachineScaleSetUpdateVMProfile{
					OsProfile: &compute.VirtualMachineScaleSetUpdateOSProfile{
						LinuxConfiguration: &compute.LinuxConfiguration{
							SSH: &compute.SSHConfiguration{
								PublicKeys: &[]compute.SSHPublicKey{
									{
										Path:    to.StringPtr("/home/cloud-user/.ssh/authorized_keys"),
										KeyData: to.StringPtr(sshKeyData),
									},
								},
							},
						},
					},
				},
			},
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func rebootInstances(ctx context.Context, logger *logrus.Entry, rg string) []error {
	var errs []error
	authorizer, err := azureclient.NewAuthorizerFromContext(ctx)
	if err != nil {
		return append(errs, err)
	}

	virtualMachineScaleSetVMs := azureclient.NewVirtualMachineScaleSetVMsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})
	virtualMachineScaleSets := azureclient.NewVirtualMachineScaleSetsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})

	logger.Debugf("listing scalesets in rg %s", rg)
	scaleSets, err := ListScaleSets(ctx, logger, rg, virtualMachineScaleSets)
	if err != nil {
		errs = append(errs, err)
	}
	for _, s := range scaleSets {
		logger.Debugf("listing vms in scale set %s", *s.Name)
		vms, err := ListVMs(ctx, logger, rg, virtualMachineScaleSetVMs, *s.Name)
		if err != nil {
			errs = append(errs, err)
		}
		for _, v := range vms {
			logger.Debugf("restarting %s", *v.Name)
			_, err := virtualMachineScaleSetVMs.Restart(ctx, rg, *s.Name, *v.ID)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}

func createScriptExtensions(ctx context.Context, logger *logrus.Entry, rg string) []error {
	var errs []error
	authorizer, err := azureclient.NewAuthorizerFromContext(ctx)
	if err != nil {
		return append(errs, err)
	}

	virtualMachineScaleSets := azureclient.NewVirtualMachineScaleSetsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})
	virtualMachineScaleSetExtensions := azureclient.NewVirtualMachineScaleSetExtensionsClient(os.Getenv("AZURE_SUBSCRIPTION_ID"), authorizer, []string{"en-us"})

	logger.Debugf("listing scale sets in rg %s", rg)
	scaleSets, err := ListScaleSets(ctx, logger, rg, virtualMachineScaleSets)
	if err != nil {
		errs = append(errs, err)
	}
	for _, s := range scaleSets {
		logger.Debugf("creating script extension for %s", *s.Name)
		_, err := virtualMachineScaleSetExtensions.CreateOrUpdate(ctx, rg, *s.Name, "test", compute.VirtualMachineScaleSetExtension{
			VirtualMachineScaleSetExtensionProperties: &compute.VirtualMachineScaleSetExtensionProperties{
				Type:     to.StringPtr("CustomScript"),
				Settings: `{"fileUris":["https://raw.githubusercontent.com/Azure-Samples/compute-automation-configurations/master/automate_nginx.sh"],"commandToExecute":"./automate_nginx.sh"}`,
			},
		})
		if err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func main() {
	opt := Options{}
	flag.StringVar(&opt.rg, "resource-group", "charlesakalugwu-prod", "The resource group provided during az openshift create")
	flag.StringVar(&opt.name, "name", "charlesakalugwu-prod", "The name of the cluster provided during az openshift create")
	flag.StringVar(&opt.location, "location", "westeurope", "The region provided during az openshift create")
	flag.StringVar(&opt.logLevel, "loglevel", "Debug", "valid values are Debug, Info, Warning, Error")

	flag.Parse()
	logrus.SetLevel(log.SanitizeLogLevel(opt.logLevel))
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	logger := logrus.WithFields(logrus.Fields{"resourceGroup": opt.rg, "cluster": opt.name, "location": opt.location})

	logger.Info("attempting to discover openshift managed resource group")

	// simulate Context with property bag
	ctx := context.Background()
	ctx = context.WithValue(ctx, api.ContextKeyClientID, os.Getenv("AZURE_CLIENT_ID"))
	ctx = context.WithValue(ctx, api.ContextKeyClientSecret, os.Getenv("AZURE_CLIENT_SECRET"))
	ctx = context.WithValue(ctx, api.ContextKeyTenantID, os.Getenv("AZURE_TENANT_ID"))

	trg := strings.Join([]string{"OS", opt.rg, opt.name, opt.location}, "_")

	mrg, err := GetManagedResourceGroup(ctx, logger, trg, opt.name, opt.location)

	if err != nil || mrg == "" {
		logger.Fatal(err)
	}

	logger.Infof("managed resource group ==> %s", mrg)
	var errs []error

	logger.Info("Updating the scale set instance count")
	errs = updateInstanceCount(ctx, logger, mrg)
	if errs != nil {
		logger.Errorf("[updateInstanceCount] %d errors %v", len(errs), errs)
	}

	logger.Info("Updating the scale set instance type")
	errs = updateInstanceType(ctx, logger, mrg)
	if errs != nil {
		logger.Errorf("[updateInstanceType] %d errors %v", len(errs), errs)
	}

	logger.Info("Updating the ssh key")
	errs = updateKey(ctx, logger, mrg)
	if errs != nil {
		logger.Errorf("[updateKey] %d errors %v", len(errs), errs)
	}

	logger.Info("Rebooting instances")
	errs = rebootInstances(ctx, logger, mrg)
	if errs != nil {
		logger.Errorf("[rebootInstances] %d errors %v", len(errs), errs)
	}

	logger.Info("Creating script extensions")
	errs = createScriptExtensions(ctx, logger, mrg)
	if errs != nil {
		logger.Errorf("[createScriptExtensions] %d errors %v", len(errs), errs)
	}
	for _, v := range errs {
		v.Error()
	}
}
