import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './index.css';

export const EditableTextField = props => {
  let content;

  if (props.isEditable) {
    content = (
      <label>
        {props.title}
        <input type="text" value={props.value} />
      </label>
    );
  } else {
    content = (
      <span>
        {props.title}: <span className="field-value">{props.value}</span>
      </span>
    );
  }

  return <div className="editable-panel-field">{content}</div>;
};

EditableTextField.propTypes = {
  title: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
};

export class EditablePanel extends Component {
  constructor() {
    super();
    this.handleToggleClick = this.handleToggleClick.bind(this);
    this.handleSaveClick = this.handleSaveClick.bind(this);
  }

  handleToggleClick(e) {
    e.preventDefault();
    this.props.toggleEditable();
  }

  handleSaveClick(e) {
    e.preventDefault();
    this.props.save();
  }

  render() {
    let className = this.props.className || '';
    className += ' editable-panel';
    let controls;

    if (this.props.isEditable) {
      className += ' is-editable';
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
              disabled={this.props.canSave}
            >
              Save
            </button>
          </p>
        </div>
      );
    }

    const ContentComponent = this.props.isEditable
      ? this.props.editableComponent
      : this.props.displayComponent;

    return (
      <div className={className}>
        <div className="editable-panel-header">
          <div className="title">{this.props.title}</div>
          {!this.props.isEditable && (
            <a className="editable-panel-edit" onClick={this.handleToggleClick}>
              Edit
            </a>
          )}
        </div>
        <div className="editable-panel-content">
          <ContentComponent />
          {controls}
        </div>
      </div>
    );
  }
}
