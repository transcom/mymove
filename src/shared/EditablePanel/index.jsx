import React, { Component } from 'react';
import PropTypes from 'prop-types';

import './index.css';

export const EditableTextField = props => {
  return (
    <div className="editable-field">
      {props.title}: {props.name}{' '}
      <em>{props.isEditable ? 'edit' : 'display'} mode</em>
    </div>
  );
};

EditableTextField.propTypes = {
  title: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
};

export class EditablePanel extends Component {
  constructor() {
    super();
    this.state = {
      isEditable: false,
    };
    this.handleClick = this.handleClick.bind(this);
    this.renderChild = this.renderChild.bind(this);
  }

  handleClick(e) {
    e.preventDefault();
    this.setState({
      isEditable: !this.state.isEditable,
    });
  }

  renderChild(child) {
    return React.cloneElement(child, { isEditable: this.state.isEditable });
  }

  render() {
    let className = 'editable-panel';
    let controls;

    if (this.state.isEditable) {
      className += ' is-editable';
      controls = (
        <div>
          <p>
            <button
              className="usa-button-secondary editable-panel-cancel"
              onClick={this.handleClick}
            >
              Cancel
            </button>
            <button className="usa-button editable-panel-save" disabled>
              Save
            </button>
          </p>
        </div>
      );
    }

    return (
      <div className={className}>
        <div className="editable-panel-header">
          <div className="title">{this.props.title}</div>
          {!this.state.isEditable && (
            <a className="editable-panel-edit" onClick={this.handleClick}>
              Edit
            </a>
          )}
        </div>
        <div className="editable-panel-content">
          {React.Children.map(this.props.children, this.renderChild)}
          {controls}
        </div>
      </div>
    );
  }
}

EditablePanel.propTypes = {
  title: PropTypes.string.isRequired,
};
