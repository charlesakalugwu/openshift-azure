### Differences between AAD Applications and Service Principal Applications

An Azure Active Directory (AAD) application is required if a user wants to 
enable logging into an OSA cluster using Azure Active Directory. This is 
enabled by configuring `AZURE_AAD_CLIENT_ID` and `AZURE_AAD_CLIENT_SECRET` in 
the env file which is sourced prior to `hack/create.sh` invocation. 

A Service Principal application on the other hand is an application that 
grants access to the OSA cluster components, enabling them to perform 
CRUD operations against the Azure Resource Provider API. The credentials 
for the Service Provider principal are specified in `AZURE_CLIENT_ID` and 
`AZURE_CLIENT_SECRET` in the env file which is sourced prior to `hack/create.sh` invocation.

Every OSA cluster requires a Service Principal application. An additional AAD 
enterprise application could be configured if the cluster operator intends to 
enable integration of the OSA cluster with AAD authentication flows. 

### Configuring your own AAD for an OpenShift Cluster on Azure

##### Using the web console
- Go to the http://portal.azure.com
- At the top of the page, type `Azure Active Directory` and click the search result to visit the Azure Active Directory Service page
- Click on `App Registrations`
- Click `New application registration`
	- Type a unique name. A good value to use here is the value that was/will be passed to `hack/create.sh` during OSA cluster creation.
	- Application type should remain `Web app /API`
	- Sign-on URL should be `https://openshift.<cluster_name>.<dns_suffix>/oauth2callback/Azure%20AD` where:
		- `<cluster_name>` is the same argument that was/will be passed to hack/create.sh  during OSA cluster creation and 
		- `<dns_suffix>` is the value of the DNS_DOMAIN from the env file that was/will be used during OSA cluster creation.
- Click `Create`
- The value of `Application ID` will be used as the value of `AZURE_AAD_CLIENT_ID` in the env file used for OSA cluster creation
- Click `Settings`
- Click `Keys`
	- In the form within the `Passwords` section, type a description, select an expiration date and click the save button. Copy the text displayed in the `Value` input field and use that as the value of `AZURE_AAD_CLIENT_SECRET` in the env file used for OSA cluster creation
- Click `Required permissions`
	- Click Add
	- Click `Select an API` 
	- Click `Microsoft Graph` and click select
	- This brings you to the `Select permissions` dialog. Select the following combinations of api/role mappings:
		- Application permissions
			- Read all groups
			- Read directory data
		- Delegated permissions
			- Read all groups
			- Sign users in
	- Click save
	- Copy the current azure portal page link from the browser and ask an admin to grant consent for your AAD application to assume those roles
	- After this you should have your env file in the root of the openshift-azure project configured with the `AZURE_AAD_CLIENT_ID` and `AZURE_AAD_CLIENT_SECRET` retrieved from the steps above.



##### Using the command line
- You can either create a new AAD app from the command line or configure your env file with a pre-created AAD application’s application id (`AZURE_AAD_CLIENT_ID`) as described in the previous step.
- To create a new AAD app from the command line:
	- Execute `./hack/aad.sh app-create <app_name> <callback_url>` where  
		- <app_name> is a unique name for the AAD application and 
		- <callback_url> is the url where a user can sign in and use our app. The form of <callback_url> should be `https://openshift.<cluster_name>.<dns_suffix>/oauth2callback/Azure%20AD` where 
			- <cluster_name> is the same argument that will be passed to hack/create.sh  during OSA cluster creation and 
			- <dns_suffix> is the value of the `DNS_DOMAIN` from the env file that was/will be used during OSA cluster creation.
	- Update your env file with the value of `AZURE_AAD_CLIENT_ID` which was output after executing the previous command
- To use a pre-created AAD application (e.g. created via the web console) simply configure its `AZURE_AAD_CLIENT_ID` variable in the env file. The following AAD-related operations will occur upon invocation of ./hack/create.sh with a configured AZURE_AAD_CLIENT_ID:
	- Calls `./hack/aad.sh app-update <aad_client_id> <callback_url>` where:
		- `<aad_client_id>` is the AAD application client id which is configured as the value of the AZURE_AAD_CLIENT_ID variable in the env file
 		- `<callback_url>` is the url where a user can sign in and use our app. The form of <callback_url> should be https://openshift.<cluster_name>.<dns_suffix>/oauth2callback/Azure%20AD where 
			- `<cluster_name>` is the same argument that will be passed to hack/create.sh  during OSA cluster creation and 
			- `<dns_suffix>` is the value of the DNS_DOMAIN from the env file that was/will be used during OSA cluster creation.
- A new client secret is generated
- The AAD app is updated with the possibly new callback url and new password
- The values for `AZURE_AAD_CLIENT_ID` and `AZURE_AAD_CLIENT_SECRET` are exported and available to the rest of the `hack/create.sh` script
- The change of password is necessary to alleviate possible man-in-the-middle attacks.
