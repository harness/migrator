---
sidebar_position: 2
---

# Prerequisites

Before using the Spinnaker to CD Next Gen migration tool, you must have the following prerequisites:

- Armory Account - Armory Platform | Spinnaker 
- Harness Account - AccountID
- Harness Org - default if not provided
- Harness Project - Project_Name
- API Key - Follow instructions at Manage API keys | [Harness Developer Hub](https://developer.harness.io/docs/platform/automation/api/add-and-manage-api-keys/#create-personal-api-keys-and-tokens) 
- Environment - Prod, Prod1, or Prod3

## Cloud Providers

Cloud providers can have their own specific prerequisites for migration. Here they are: 

### AWS

If you are using the **spinnaker-aws-drone-plugin** for any stages (i.e. Shrinkcluster, Scaledowncluster, Invokelambda) then ensure that there is an AWS secret created at the project level within Harness with the following names:

```
AWS_ACCESS_KEY_ID
AWS_SECRET_ACCESS_KEY
```

For information on how to create secrets with Harness, go to [AWS Secrets Manager](https://developer.harness.io/docs/platform/secrets/secrets-management/add-an-aws-secret-manager)

### GCP

Regrettably, Harness lacks native support for triggering pipelines directly via incoming Pub/Sub messages. However, we can work around this limitation by using a custom webhook trigger. 

For detailed information about custom webhook triggers, go to [Using Custom Triggers](https://developer.harness.io/docs/platform/triggers/trigger-deployments-using-custom-triggers/)

In this approach, you would need to provision a new GCP Cloud Function that triggers from Pub/Sub topic messages and invokes the custom webhook URL. Below is an example of such a function in Node.js.

```jsx
const fetch = require('node-fetch');

exports.pubsubListener = (event, context) => {
  const pubsubMessage = event.data;
  const message = Buffer.from(pubsubMessage, 'base64').toString();

  console.log(`Received message: ${message}`);

  // Define the webhook URL and payload
  const webhookUrl = '<https://app.harness.io/gateway/pipeline/api/webhook/custom/Wlke4SU6TAeKwDERALcptQ/v3?accountIdentifier=YOUR_ACCOUNT_ID&orgIdentifier=YOUR_ORG_ID&projectIdentifier=YOUR_PROJECT_ID&pipelineIdentifier=YOUR_PIPELINE_ID&triggerIdentifier=webhook_trigger>';
  const payload = {
    sample_key: 'sample_value'
  };

  // Make a POST request to the webhook URL
  fetch(webhookUrl, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  })
  .then(response => {
    if (!response.ok) {
      throw new Error('Failed to trigger webhook');
    }
    console.log('Webhook triggered successfully');
    return response.json();
  })
  .catch(error => {
    console.error('Error triggering webhook:', error);
  });
};
```
