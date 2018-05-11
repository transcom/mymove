import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import './index.css';

export const PanelField = props => {
  return (
    <div className="panel-field">
      <span className="field-title">{props.title}</span>
      <span className="field-value">{props.value}</span>
    </div>
  );
};

export class EditablePanel extends Component {
  handleToggleClick = e => {
    e.preventDefault();
    this.props.toggleEditable();
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
  toggleEditable: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
};
