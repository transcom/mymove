import { get } from 'lodash';
import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { editablePanelify } from 'shared/EditablePanel';
import { createOrUpdateServiceAgent } from './ducks';

import { PanelSwaggerField } from 'shared/EditablePanel';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const ServiceAgentDisplay = props => {
  let serviceAgent = props.initialValues || {};

  const fieldProps = {
    schema: props.schema,
    values: serviceAgent,
  };
  return (
    <Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField
          fieldName="point_of_contact"
          nullWarning
          {...fieldProps}
        />
        <PanelSwaggerField fieldName="email" nullWarning {...fieldProps} />
        <PanelSwaggerField
          fieldName="phone_number"
          nullWarning
          {...fieldProps}
        />
      </div>
    </Fragment>
  );
};

const ServiceAgentEdit = props => {
  const schema = props.schema;
  return (
    <Fragment>
      <div className="editable-panel-column">
        <SwaggerField fieldName="point_of_contact" swagger={schema} required />
        <SwaggerField fieldName="email" swagger={schema} required />
        <SwaggerField fieldName="phone_number" swagger={schema} required />
      </div>
    </Fragment>
  );
};

const formName = 'service_agents_panel';

let ServiceAgentPanel = editablePanelify(ServiceAgentDisplay, ServiceAgentEdit);
ServiceAgentPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(ServiceAgentPanel);

const ServiceAgents = props => {
  const { schema } = props;

  if (!schema) {
    return <LoadingPlaceholder />;
  }

  let originAgent = { role: 'ORIGIN' };
  let destinationAgent = { role: 'DESTINATION' };

  // This returns the last value. Currently there should only be one record for each role.
  props.serviceAgents.forEach(agent => {
    if (agent.role === 'ORIGIN') {
      originAgent = agent;
    } else if (agent.role === 'DESTINATION') {
      destinationAgent = agent;
    } else {
      console.error('Unknown Agent Role: ', agent);
    }
  });

  return (
    <Fragment>
      <ServiceAgentPanel
        form="origin_service_agent"
        title="Origin Service Agent"
        update={props.update}
        agentRole="ORIGIN"
        schema={props.schema}
        initialValues={originAgent}
        getUpdateArgs={props.getOriginUpdateArgs}
      />

      <ServiceAgentPanel
        form="destination_service_agent"
        title="Destination Service Agent"
        update={props.update}
        agentRole="DESTINATION"
        schema={props.schema}
        initialValues={destinationAgent}
        getUpdateArgs={props.getDestinationUpdateArgs}
      />
    </Fragment>
  );
};

function mapStateToProps(state, props) {
  let originFormValues = getFormValues('origin_service_agent')(state);
  let destFormValues = getFormValues('destination_service_agent')(state);

  return {
    // reduxForm
    schema: get(state, 'swagger.spec.definitions.ServiceAgent', null),

    hasError: false,
    errorMessage: state.tsp.error,
    isUpdating: false,

    // editablePanel
    getOriginUpdateArgs: function() {
      return [get(props, 'shipment.id'), originFormValues];
    },
    getDestinationUpdateArgs: function() {
      return [get(props, 'shipment.id'), destFormValues];
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
