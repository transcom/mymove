import yaml from 'js-yaml';

//this is kind of clunky, it might be better to just convert this to json
export const getUiSchema = () => yaml.safeLoad(uiYaml);

const uiYaml = `
order:
  - feedback
groups:
  feedback:
    title: Customer Feedback
    fields:
      - description
      - reporter_name
      - due_date
      - telephone
      - annoyance_level
      - some_uuid
`;
