import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { get } from 'lodash';

import { formatCents } from 'shared/formatters';
import { formatDate } from 'utils/formatters';
import Alert from 'shared/Alert';

import './index.css';

const DefaultHeader = ({ title, isEditable, editEnabled, handleEditClick }) => (
  <div className="editable-panel-header">
    <h4>{title}</h4>
    {!isEditable && editEnabled && (
      <a data-testid="edit-link" className="usa-link editable-panel-edit" onClick={handleEditClick}>
        Edit
      </a>
    )}
  </div>
);

export const RowBasedHeader = ({ title, isEditable, editEnabled, handleEditClick }) => (
  <div className="grid-row">
    <div className="grid-col-10">
      <h1>{title}</h1>
    </div>
    {!isEditable && editEnabled && (
      <div className="grid-col-2 text-right">
        <p>
          <a data-testid="edit-link" className="usa-link" onClick={handleEditClick}>
            Edit
          </a>
        </p>
      </div>
    )}
  </div>
);

export const PanelField = (props) => {
  const { title, value, required } = props;
  const classes = classNames('panel-field', props.className);
  let component = (
    <div className={classes}>
      <span className="field-title">{title}</span>
      <span className="field-value">{value || props.children}</span>
    </div>
  );

  if (required && !(value || props.children)) {
    component = (
      <div className={'missing ' + classes}>
        <span className="field-title">{title}</span>
        <span className="field-value">missing</span>
      </div>
    );
  }

  return component;
};
PanelField.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.string,
  children: PropTypes.node,
  className: PropTypes.string,
  required: PropTypes.bool,
};

export const SwaggerValue = (props) => {
  const { fieldName, schema, values } = props;
  let swaggerProps = {};
  if (schema.properties) {
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
  return <React.Fragment>{value || null}</React.Fragment>;
};
SwaggerValue.propTypes = {
  fieldName: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  values: PropTypes.object,
};

export const PanelSwaggerField = (props) => {
  const { fieldName, className, required, schema, values } = props;
  const title = props.title || get(schema, `properties.${fieldName}.title`, fieldName);
  const classes = classNames(fieldName, className);
  let component = (
    <PanelField title={title} className={classes}>
      <SwaggerValue {...props} />
      {props.children}
    </PanelField>
  );

  if (required && !values[fieldName]) {
    component = (
      <PanelField title={title} className={classes} required>
        {props.children}
      </PanelField>
    );
  }

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
  handleEditClick = (e) => {
    e.preventDefault();
    this.props.onEdit();
  };

  handleCancelClick = (e) => {
    e.preventDefault();
    this.props.onCancel();
  };

  handleSaveClick = (e) => {
    e.preventDefault();
    this.props.onSave();
  };

  render() {
    let controls;

    if (this.props.isEditable) {
      controls = (
        <div>
          <p>
            <button className="usa-button usa-button--secondary editable-panel-cancel" onClick={this.handleCancelClick}>
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

    const { title, isEditable, editEnabled, HeaderComponent } = this.props;

    return (
      <div className={classes}>
        <HeaderComponent
          title={title}
          isEditable={isEditable}
          editEnabled={editEnabled}
          handleEditClick={this.handleEditClick}
        />
        <div className="editable-panel-content">
          {this.props.children}
          {controls}
        </div>
      </div>
    );
  }
}

// Convenience function for creating an editable panel given a display component and an edit component
export function editablePanelify(DisplayComponent, EditComponent, editEnabled = true, HeaderComponent = DefaultHeader) {
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
      this.setState((prevState) => ({ isEditable: !prevState.isEditable }));
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
            HeaderComponent={HeaderComponent}
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
