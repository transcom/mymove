import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { get } from 'lodash';

import './index.css';

export const PanelField = props => {
  const { title, value } = props;
  const classes = classNames('panel-field', props.className);
  return (
    <div className={classes}>
      <span className="field-title">{title}</span>
      <span className="field-value">{value || props.children}</span>
    </div>
  );
};
PanelField.propTypes = {
  title: PropTypes.string.isRequired,
  value: PropTypes.string,
  children: PropTypes.node,
  className: PropTypes.string,
};

export const SwaggerValue = props => {
  const { fieldName, schema, values } = props;
  /* eslint-disable security/detect-object-injection */
  const swaggerProps = schema.properties[fieldName];

  let value = values[fieldName];
  if (swaggerProps.enum) {
    value = swaggerProps['x-display-value'][value];
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
  const { fieldName, schema } = props;
  const title =
    props.title || get(schema, `properties.${fieldName}.title`, fieldName);

  return (
    <PanelField title={title}>
      <SwaggerValue {...props} />
    </PanelField>
  );
};
PanelSwaggerField.propTypes = {
  fieldName: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  values: PropTypes.object.isRequired,
  title: PropTypes.string,
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
            <button
              className="usa-button-secondary editable-panel-cancel"
              onClick={this.handleCancelClick}
            >
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
          {!this.props.isEditable && (
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

EditablePanel.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node.isRequired,
  isEditable: PropTypes.bool.isRequired,
  isValid: PropTypes.bool.isRequired,
  onCancel: PropTypes.func.isRequired,
  onEdit: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
};
