import React, { Component } from 'react';
import PropTypes from 'prop-types';

class DeleteSitRequest extends Component {
  render() {
    return (
      <div className="usa-alert usa-alert-warning sit-delete-warning">
        <div className="sit-delete-buttons">
          <button className="usa-button">Yes, Delete</button>
          &nbsp;&nbsp;
          <a className="sit-delete-cancel" onClick={this.props.onClose}>
            No, do not delete
          </a>
        </div>
        <div className="usa-alert-body">
          <h3 className="usa-alert-heading">Delete this SIT request?</h3>
          <p className="usa-alert-text">This action cannot be undone.</p>
        </div>
      </div>
    );
  }
}

DeleteSitRequest.propTypes = {
  storageInTransit: PropTypes.object.isRequired,
  onClose: PropTypes.func.isRequired,
};

export default DeleteSitRequest;
