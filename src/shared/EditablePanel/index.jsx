import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';
import { get } from 'lodash';

import './index.css';

export const PanelField = props => {
  const { fieldName, schema, values } = props;
  const title = get(schema, `properties.${fieldName}.title`, '');
  const value = values[fieldName];

  return (
    <div className="panel-field">
      <span className="field-title">{title}</span>
      <span className="field-value">{value}</span>
    </div>
  );
};
PanelField.propTypes = {
  fieldName: PropTypes.string.isRequired,
  schema: PropTypes.object.isRequired,
  values: PropTypes.object,
};

export class EditablePanel extends Component {
  handleToggleClick = e => {
    e.preventDefault();
    this.props.onToggle();
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
              onClick={this.handleToggleClick}
            >
              Cancel
            </button>
            <button
              className="usa-button editable-panel-save"
              onClick={this.handleSaveClick}
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
            <a className="editable-panel-edit" onClick={this.handleToggleClick}>
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
  onToggle: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
};
