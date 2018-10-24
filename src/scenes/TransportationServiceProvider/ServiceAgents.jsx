import React, { Component, Fragment } from 'react';
import { get } from 'lodash';
import classNames from 'classnames';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm, getFormValues } from 'redux-form';
import PropTypes from 'prop-types';
import Alert from 'shared/Alert';

import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import { createOrUpdateServiceAgent } from './ducks';

import { PanelSwaggerField } from 'shared/EditablePanel';

import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

// TODO: Refactor when we switch to using a wizard
// Editable panel specific to Assign Service Agents. Due to not using a wizard to assign the service agents this
// panel has highly specific behavior (opening the edit view via clicking on Assign Service Agents button)
export class ServiceAgentEditablePanel extends Component {
  handleEditClick = e => {
    e.preventDefault();
    this.props.onEdit(true);
  };

  handleCancelClick = e => {
    e.preventDefault();
    this.props.onCancel();
  };

  handleSaveClick = e => {
    e.preventDefault();
    this.props.onSave();
  };

  render() {
    let controls;

    if (this.props.isEditable) {
      controls = (
        <div>
          <p>
            <button className="usa-button-secondary editable-panel-cancel" onClick={this.handleCancelClick}>
              Cancel
            </button>
            <button
              className="usa-button editable-panel-save"
              onClick={this.handleSaveClick}
              disabled={!this.props.isValid}
            >
              Save
            </button>
          </p>
        </div>
      );
    }

    const classes = classNames(
      'editable-panel',
      {
        'is-editable': this.props.isEditable,
      },
      this.props.className,
    );

    return (
      <div className={classes}>
        <div className="editable-panel-header">
          <div className="title">{this.props.title}</div>
          {!this.props.isEditable &&
            this.props.editEnabled && (
              <a className="editable-panel-edit" onClick={this.handleEditClick}>
                Edit
              </a>
            )}
        </div>
        <div className="editable-panel-content">
          {this.props.children}
          {controls}
        </div>
      </div>
    );
  }
}

// TODO: Refactor when we switch to using a wizard
// Editable panel specific to Assign Servivce Agents. Due to not using a wizard to assign the service agents this
// panel has highly specific behavior (opening the edit view via clicking on Assign Service Agents button)
export function ServiceAgentEditablePanelify(DisplayComponent, EditComponent, editEnabled = true) {
  const Wrapper = class extends Component {
    state = {
      isEditable: false,
    };

    componentDidUpdate = (prevProps, prevState) => {
      if (!prevProps.editOriginServiceAgent && this.props.editOriginServiceAgent) {
        this.setIsEditable(true);
      }
    };

    save = () => {
      let isValid = this.props.valid;
      if (isValid) {
        let args = this.props.getUpdateArgs();
        this.props.update(...args);
        this.setIsEditable(false);
      }
    };

    cancel = () => {
      this.props.reset();
      this.setIsEditable(false);
      if (this.props.title === 'Origin Service Agent') {
        this.props.setEditServiceAgent(false);
      }
    };

    setIsEditable = isEditable => this.setState({ isEditable });

    render() {
      const isEditable = (editEnabled && (this.state.isEditable || this.props.isUpdating)) || false;
      const Content = isEditable ? EditComponent : DisplayComponent;

      return (
        <React.Fragment>
          {this.props.hasError && (
            <Alert type="error" heading="An error occurred">
              There was an error: <em>{this.props.errorMessage}</em>.
            </Alert>
          )}
          <ServiceAgentEditablePanel
            title={this.props.title}
            className={this.props.className}
            onSave={this.save}
            onEdit={this.setIsEditable}
            onCancel={this.cancel}
            isEditable={isEditable}
            editEnabled={editEnabled}
            isValid={this.props.valid}
          >
            <Content {...this.props} />
          </ServiceAgentEditablePanel>
        </React.Fragment>
      );
    }
  };

  Wrapper.propTypes = {
    update: PropTypes.func.isRequired,
    title: PropTypes.string.isRequired,
    isUpdating: PropTypes.bool,
  };

  return Wrapper;
}

ServiceAgentEditablePanel.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  isEditable: PropTypes.bool.isRequired,
  editEnabled: PropTypes.bool,
  isValid: PropTypes.bool,
  onCancel: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
};

ServiceAgentEditablePanel.defaultProps = {
  editEnabled: true,
};

const ServiceAgentDisplay = props => {
  let serviceAgent = props.initialValues || {};

  const fieldProps = {
    schema: props.schema,
    values: serviceAgent,
  };
  return (
    <Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField fieldName="company" required {...fieldProps} />
        <PanelSwaggerField fieldName="email" required {...fieldProps} />
        <PanelSwaggerField fieldName="phone_number" required {...fieldProps} />
      </div>
    </Fragment>
  );
};

const ServiceAgentEdit = props => {
  const schema = props.schema;
  return (
    <Fragment>
      <div className="editable-panel-column">
        <SwaggerField fieldName="company" swagger={schema} required />
        <SwaggerField fieldName="email" swagger={schema} required />
        <SwaggerField fieldName="phone_number" swagger={schema} required />
      </div>
    </Fragment>
  );
};

const formName = 'service_agents_panel';

let ServiceAgentPanel = ServiceAgentEditablePanelify(ServiceAgentDisplay, ServiceAgentEdit);

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
        editOriginServiceAgent={props.editOriginServiceAgent}
        setEditServiceAgent={props.setEditServiceAgent}
        update={props.update}
        schema={props.schema}
        initialValues={props.initialValues.ORIGIN}
        getUpdateArgs={props.getOriginUpdateArgs}
      />

      <ServiceAgentPanel
        form="destination_service_agent"
        title="Destination Service Agent"
        update={props.update}
        schema={props.schema}
        initialValues={props.initialValues.DESTINATION}
        getUpdateArgs={props.getDestinationUpdateArgs}
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
    schema: get(state, 'swaggerPublic.spec.definitions.ServiceAgent', null),
    initialValues,

    hasError: false,
    errorMessage: get(state, 'tsp.error', {}),
    isUpdating: false,

    // editablePanel
    getOriginUpdateArgs: function() {
      return [get(props, 'shipment.id'), Object.assign({}, originFormValues, { role: 'ORIGIN' })];
    },
    getDestinationUpdateArgs: function() {
      return [get(props, 'shipment.id'), Object.assign({}, destFormValues, { role: 'DESTINATION' })];
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
