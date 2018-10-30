import { get } from 'lodash';
import React, { Component, Fragment } from 'react';
import classNames from 'classnames';
import { reduxForm, FormSection } from 'redux-form';
import PropTypes from 'prop-types';
import Alert from 'shared/Alert';

import { ServiceAgentDisplay, ServiceAgentEdit } from './ServiceAgentViews';

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
      if (!prevProps.editServiceAgent && this.props.editServiceAgent) {
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

const TSPDisplay = props => {
  const originSAProps = {
    schema: props.saSchema,
    values: props.origin_service_agent,
  };

  const destinationSAProps = {
    schema: props.saSchema,
    values: props.destination_service_agent,
  };
  // const TSPprops = {
  //   schame: props.tspSchema,
  //   values: TSP
  // }
  return (
    <Fragment>
      <ServiceAgentDisplay serviceAgentProps={originSAProps} role="Origin" />
      <ServiceAgentDisplay serviceAgentProps={destinationSAProps} role="Destination" />
    </Fragment>
  );
};

const TSPEdit = props => {
  const saSchema = props.saSchema;
  const originValues = get(props, 'formValues.ORIGIN', {});
  const destinationValues = get(props, 'formValues.DESTINATION', {});

  return (
    <Fragment>
      <FormSection name="origin_service_agent">
        <ServiceAgentEdit
          serviceAgentProps={{
            swagger: saSchema,
            values: originValues,
          }}
          role="Origin"
        />
      </FormSection>
      <FormSection name="destination_service_agent">
        <ServiceAgentEdit serviceAgentProps={{ swagger: saSchema, values: destinationValues }} role="Destination" />
      </FormSection>
    </Fragment>
  );
};

let TspPanel = ServiceAgentEditablePanelify(TSPDisplay, TSPEdit);

TspPanel = reduxForm({
  // formName passed in implicitly by container
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(TspPanel);

export default TspPanel;
