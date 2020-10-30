// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import './UploadsTable.css';

import bytes from 'bytes';
import moment from 'moment';
import { UPLOAD_SCAN_STATUS } from 'shared/constants';

export class UploadsTable extends Component {
  getUploadUrl = (upload) => {
    if (upload.status === UPLOAD_SCAN_STATUS.INFECTED) {
      return (
        <>
          <Link to="/infected-upload" className="usa-link">
            {upload.filename}
          </Link>
        </>
      );
    } else if (upload.status === UPLOAD_SCAN_STATUS.PROCESSING) {
      return (
        <>
          <Link to="/processing-upload" className="usa-link">
            {upload.filename}
          </Link>
        </>
      );
    }
    return (
      <>
        <a href={upload.url} target="_blank" rel="noopener noreferrer" className="usa-link">
          {upload.filename}
        </a>
      </>
    );
  };

  render() {
    return (
      <table>
        <thead>
          <tr>
            <th>Name</th>
            <th>Uploaded</th>
            <th>Size</th>
            <th>Delete</th>
          </tr>
        </thead>
        <tbody>
          {this.props.uploads.map((upload) => (
            <tr key={upload.id} className="vertical-align text-top">
              <td className="maxw-card" style={{ overflowWrap: 'break-word', wordWrap: 'break-word' }}>
                {this.getUploadUrl(upload)}
              </td>
              <td>{moment(upload.created_at).format('LLL')}</td>
              <td>{bytes(upload.bytes)}</td>
              <td>
                <a href="" onClick={(e) => this.props.onDelete(e, upload.id)} className="usa-link">
                  Delete
                </a>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    );
  }
}

UploadsTable.propTypes = {
  uploads: PropTypes.array.isRequired,
  onDelete: PropTypes.func,
};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(UploadsTable);
