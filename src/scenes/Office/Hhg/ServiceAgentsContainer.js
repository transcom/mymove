import { connect } from 'react-redux';
import { get, isEmpty } from 'lodash';
import { getFormValues } from 'redux-form';
import { bindActionCreators } from 'redux';

import TspServiceAgents from 'shared/TspPanel/TspServiceAgents';
import { updateServiceAgentsForShipment } from 'shared/Entities/modules/serviceAgents';

import { getPublicSwaggerDefinition } from 'shared/Swagger/selectors';
import { selectTspById } from 'shared/Entities/modules/transportationServiceProviders';

function mapStateToProps(state, props) {
  let serviceAgents = props.serviceAgents;
  let title = props.title;
  let initialValues = {};
  let form = 'tsp_service_agents';
  let formValues = getFormValues(form)(state);
  // This returns the last value. Currently there should only be one record for each role.
  serviceAgents.forEach(agent => (initialValues[agent.role] = agent));
  const { ORIGIN: origin_service_agent = {}, DESTINATION: destination_service_agent = {} } = initialValues;

  return {
    // reduxForm
    saSchema: getPublicSwaggerDefinition(state, 'ServiceAgent', null),
    transportationServiceProvider: selectTspById(state, props.transportationServiceProviderId),
    formValues,
    initialValues: {
      origin_service_agent,
      destination_service_agent,
    },
    origin_service_agent,
    destination_service_agent,
    title,
    form,

    hasError: false,
    errorMessage: get(state, 'tsp.error', {}),
    isUpdating: false,

    // editablePanel
    getUpdateArgs: function() {
      const params = {
        origin_service_agent: { ...formValues.origin_service_agent, role: 'ORIGIN' },
      };
      // Avoid sending empty request for destination service agent
      if (!isEmpty(formValues.destination_service_agent)) {
        params.destination_service_agent = { ...formValues.destination_service_agent, role: 'DESTINATION' };
      }

      return [get(props, 'shipment.id'), params];
    },
  };
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: updateServiceAgentsForShipment,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(TspServiceAgents);
