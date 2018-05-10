import React, { Component } from 'react';
import PropTypes from 'prop-types';
import classNames from 'classnames';

import './index.css';

export class EditablePanel extends Component {
  handleToggleClick = e => {
    e.preventDefault();
    this.props.toggleEditable();
  };

  // TODO: get save button working as more than just an editable toggle
  handleSaveClick = e => {
    e.preventDefault();
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
              onClick={this.handleToggleClick}
              disabled={this.props.canSave}
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

    const contentFunc = this.props.isEditable
      ? this.props.editableContent
      : this.props.displayContent;

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
          {contentFunc()}
          {controls}
        </div>
      </div>
    );
  }
}

EditablePanel.propTypes = {
  isEditable: PropTypes.bool.isRequired,
  editableContent: PropTypes.func.isRequired,
  displayContent: PropTypes.func.isRequired,
};
