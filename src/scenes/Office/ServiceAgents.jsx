import { get } from 'lodash';
import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { editablePanelify, PanelSwaggerField } from 'shared/EditablePanel';
import { no_op } from 'shared/utils';

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

const ServiceAgentEdit = 'uneditable';
const formName = 'service_agents_panel';
const editEnabled = false;
let ServiceAgentPanel = editablePanelify(
  ServiceAgentDisplay,
  ServiceAgentEdit,
  editEnabled,
);
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

  return (
    <Fragment>
      <ServiceAgentPanel
        form="origin_service_agent"
        title="Origin Service Agent"
        update={props.update}
        schema={props.schema}
        initialValues={props.initialValues.ORIGIN}
        getUpdateArgs={props.getOriginUpdateArgs}
        editEnabled={props.editEnabled}
      />

      <ServiceAgentPanel
        form="destination_service_agent"
        title="Destination Service Agent"
        update={props.update}
        schema={props.schema}
        initialValues={props.initialValues.DESTINATION}
        getUpdateArgs={props.getDestinationUpdateArgs}
        editEnabled={props.editEnabled}
      />
    </Fragment>
  );
};

function mapStateToProps(state, props) {
  let originFormValues = getFormValues('origin_service_agent')(state);
  let destFormValues = getFormValues('destination_service_agent')(state);
  let serviceAgents = props.serviceAgents;
  let initialValues = {};
  // This returns the last value. Currently there should only be one record for each role.
  serviceAgents.forEach(agent => (initialValues[agent.role] = agent));

  return {
    // reduxForm
    schema: get(state, 'swagger.spec.definitions.ServiceAgent', null),
    initialValues,

    hasError: false,
    errorMessage: get(state, 'shipment.error', {}),
    isUpdating: false,

    // editablePanel
    getOriginUpdateArgs: function() {
      return [
        get(props, 'shipment.id'),
        Object.assign({}, originFormValues, { role: 'ORIGIN' }),
      ];
    },
    getDestinationUpdateArgs: function() {
      return [
        get(props, 'shipment.id'),
        Object.assign({}, destFormValues, { role: 'DESTINATION' }),
      ];
    },
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(ServiceAgents);
