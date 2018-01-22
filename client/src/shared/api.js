import jsYaml from 'js-yaml';
import Swagger from 'swagger-client';

const getApi = () =>
  fetch('api/v1/swagger.yaml')
    .then(response => response.text())
    .then(yaml => jsYaml.safeLoad(yaml))
    .then(json => Swagger({ spec: json }));

export const IssuesIndex = () =>
  getApi().then(jx => jx.apis.default.indexIssues());
