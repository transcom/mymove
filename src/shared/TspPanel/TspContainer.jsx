import { connect } from 'react-redux';
import { get } from 'lodash';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';

import { createOrUpdateServiceAgent } from 'scenes/TransportationServiceProvider/ducks';
import ServiceAgents from './ServiceAgents';

import { getPublicSwaggerDefinition } from 'shared/Swagger/selectors';

function mapStateToProps(state, props) {
  let serviceAgents = props.serviceAgents;
  let initialValues = {};
  let form = 'tsp_service_agents';
  let formValues = getFormValues(form)(state);
  // This returns the last value. Currently there should only be one record for each role.
  serviceAgents.forEach(agent => (initialValues[agent.role] = agent));
  const { ORIGIN: origin_service_agent = {}, DESTINATION: destination_service_agent = {} } = initialValues;

  return {
    // reduxForm
    saSchema: getPublicSwaggerDefinition(state, 'ServiceAgent', null),
    tspSchema: getPublicSwaggerDefinition(state, 'TransportationServiceProvider', null),
    formValues,
    initialValues: {
      origin_service_agent,
      destination_service_agent,
    },
    origin_service_agent,
    destination_service_agent,
    title: 'TSP & Servicing Agents',
    form,

    hasError: false,
    errorMessage: get(state, 'tsp.error', {}),
    isUpdating: false,

    // editablePanel
    getUpdateArgs: function() {
      const params = {
        origin_service_agent: formValues.origin_service_agent,
        destination_service_agent: formValues.destination_service_agent,
      };
      return [get(props, 'shipment.id'), params];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: createOrUpdateServiceAgent,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceAgents);
