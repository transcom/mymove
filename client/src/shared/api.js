import Swagger from 'swagger-client';

export const IssuesIndex = () =>
  Swagger('api/v1/swagger.yaml').then(client => {
    return client.apis.default.indexIssues();
  });

export const CreateIndex = issueBody =>
  Swagger('api/v1/swagger.yaml').then(client => {
    return client.apis.default.createIssue({ body: issueBody });
  });
