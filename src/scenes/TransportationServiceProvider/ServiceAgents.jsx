import { get } from 'lodash';
import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { editablePanelify } from 'shared/EditablePanel';
import { createServiceAgent } from './ducks';

import { PanelSwaggerField } from 'shared/EditablePanel';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const ServiceAgentDisplay = props => {
  let serviceAgent = props.initialValues;
  let role = props.agentRole;

  const fieldProps = {
    schema: props.schema,
    values: serviceAgent,
  };
  return (
    <Fragment>
      <span>{role}</span>
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

  const fakeServiceAgent = {
    point_of_contact: 'some guy',
    email: 'james@baxter.com',
    phone_number: null,
    role: 'ORIGIN',
  };

  if (!schema) {
    return <LoadingPlaceholder />;
  }

  return (
    <Fragment>
      <h1>Service Agents!</h1>

      <ServiceAgentPanel
        title="Origin Service Agent"
        update={props.update}
        formIsValid={props.formIsValid}
        agentRole="ORIGIN"
        schema={props.schema}
        initialValues={fakeServiceAgent}
        getUpdateArgs={props.getUpdateArgs}
      />
    </Fragment>
  );
};

function mapStateToProps(state, props) {
  let formValues = getFormValues(formName)(state);

  const serviceAgents = get(state, 'tsp.serviceAgents', []);

  return {
    // reduxForm
    formValues: formValues,
    initialValues: {},

    schema: get(state, 'swagger.spec.definitions.ServiceAgent', null),

    hasError: false,
    errorMessage: state.tsp.error,
    isUpdating: false,

    serviceAgents,

    // editablePanel
    getUpdateArgs: function() {
      return [get(props, 'shipment.id'), formValues];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: createServiceAgent,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceAgents);
