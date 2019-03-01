import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { get } from 'lodash';

import { formatCents, formatDate } from 'shared/formatters';
import Alert from 'shared/Alert';

import './index.css';

export const PanelField = props => {
  const { title, value, required } = props;
  const classes = classNames('panel-field', props.className);
  let component = (
    <div className={classes}>
      <span className="field-title">{title}</span>
      <span className="field-value">{value || props.children}</span>
    </div>
  );

  /* eslint-disable security/detect-object-injection */
  if (required && !(value || props.children)) {
    component = (
      <div className={'missing ' + classes}>
        <span className="field-title">{title}</span>
        <span className="field-value">missing</span>
      </div>
    );
  }
  /* eslint-enable security/detect-object-injection */

  return component;
};
PanelField.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.string,
  children: PropTypes.node,
  className: PropTypes.string,
  required: PropTypes.bool,
};

export const SwaggerValue = props => {
  const { fieldName, schema, values } = props;
  let swaggerProps = {};
  if (schema.properties) {
    /* eslint-disable security/detect-object-injection */
    swaggerProps = schema.properties[fieldName];
  }
  let value = values[fieldName] || '';
  if (swaggerProps.enum) {
    value = swaggerProps['x-display-value'][value];
  }
  if (swaggerProps.format === 'cents') {
    value = formatCents(value);
  }
  if (swaggerProps.format === 'date') {
    value = formatDate(value);
  }
  if (value && swaggerProps['x-formatting'] === 'weight') {
    value = value.toLocaleString() + ' lbs';
  }
  /* eslint-enable security/detect-object-injection */
  return <React.Fragment>{value || null}</React.Fragment>;
};
SwaggerValue.propTypes = {
  fieldName: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  values: PropTypes.object,
};

export const PanelSwaggerField = props => {
  const { fieldName, className, required, schema, values } = props;
  const title = props.title || get(schema, `properties.${fieldName}.title`, fieldName);
  const classes = classNames(fieldName, className);
  let component = (
    <PanelField title={title} className={classes}>
      <SwaggerValue {...props} />
      {props.children}
    </PanelField>
  );

  /* eslint-disable security/detect-object-injection */
  if (required && !values[fieldName]) {
    component = (
      <PanelField title={title} className={classes} required>
        {props.children}
      </PanelField>
    );
  }
  /* eslint-enable security/detect-object-injection */

  return component;
};
PanelSwaggerField.propTypes = {
  fieldName: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  values: PropTypes.object.isRequired,
  title: PropTypes.string,
  children: PropTypes.node,
  required: PropTypes.bool,
  className: PropTypes.string,
};

export class EditablePanel extends Component {
  handleEditClick = e => {
    e.preventDefault();
    this.props.onEdit();
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
      this.props.title.toLowerCase(),
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

// Convenience function for creating an editable panel given a display component and an edit component
export function editablePanelify(DisplayComponent, EditComponent, editEnabled = true) {
  const Wrapper = class extends Component {
    state = {
      isEditable: false,
    };

    save = () => {
      let isValid = this.props.valid;
      if (isValid) {
        let args = this.props.getUpdateArgs();
        this.props.update(...args);
        this.toggleEdit();
      }
    };

    toggleEdit = () => {
      this.setState({ isEditable: !this.state.isEditable });
    };

    cancel = () => {
      this.props.reset();
      this.toggleEdit();
    };

    render() {
      const isEditable = (editEnabled && (this.state.isEditable || this.props.isUpdating)) || false;
      const Content = isEditable ? EditComponent : DisplayComponent;

      return (
        <React.Fragment>
          {this.props.hasError && (
            <Alert type="error" heading="An error occurred">
              <em>{this.props.errorMessage}</em>
            </Alert>
          )}
          <EditablePanel
            title={this.props.title}
            className={this.props.className}
            onSave={this.save}
            onEdit={this.toggleEdit}
            onCancel={this.cancel}
            isEditable={isEditable}
            editEnabled={editEnabled}
            isValid={this.props.valid}
          >
            <Content {...this.props} />
          </EditablePanel>
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

EditablePanel.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  isEditable: PropTypes.bool.isRequired,
  editEnabled: PropTypes.bool,
  isValid: PropTypes.bool,
  onCancel: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
};

EditablePanel.defaultProps = {
  editEnabled: true,
};
