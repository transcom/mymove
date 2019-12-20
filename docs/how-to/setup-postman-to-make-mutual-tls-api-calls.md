# How to Setup Postman to make Mutual TLS API Calls

If you are planning to use [Postman](https://www.getpostman.com/) for testing the api you will need to make the following changes to support Mutual TLS.

## General Postman Settings

Open the general settings panel by clicking the wrench icon in the upper left corner

![Postman Settings Menu Upper Right Corner](../images/postman_settings_menu.png)

Under the _General_ tab turn off **SSL certificate verification**

![Postman SSL certification verification switch](../images/postman_ssl_verification.png)

Switch to the **Certificates** tab and add the development certificate with the following settings:

* **Host** `primelocal`
* **Port** `9443`
* **CRT File** `config/tls/devlocal-mtls.cer`
* **KEY File** `config/tls/devlocal-mtls.key`

![Postman client cert settings](../images/postman_client_cert.png)

## Postman Environment settings

You will need to configure the base url for development or other environment you plan to connect to. Click on the gear icon near the environment pull down in the upper right of the application.

![Postman open environment dialog](../images/postman_environment.png)

This will open the _Manage Environments_ dialog. Select **Add** in the lower right corner

![Postman environment dialog](../images/postman_manage_environment_dialog.png)

Fill in the following details in the add new dialog and click **Add**

* **Variable** `baseUrl`
* **Initial Value** `https://primelocal:9443/prime/v1`
* **Current Value** `https://primelocal:9443/prime/v1`

![Postman environment add dialog](../images/postman_manage_environment_add.png)

Once you have added this environment and closed the dialog select the new environment from the pull down.

![Postman select environment](../images/postman_set_environment.png)
