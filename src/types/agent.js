import { string, shape } from 'prop-types';

export const AgentShape = shape({
  agentType: string.isRequired,
  firstName: string.isRequired,
  lastName: string,
  phone: string,
  email: string,
});

export default {
  AgentShape,
};
