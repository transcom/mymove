// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';

import './UploadsTable.css';

import bytes from 'bytes';
import moment from 'moment';
import { UPLOAD_SCAN_STATUS } from 'shared/constants';

export class UploadsTable extends Component {
  getUploadUrl = upload => {
    if (upload.status === UPLOAD_SCAN_STATUS.INFECTED) {
      return `/infected-upload`;
    } else if (upload.status === UPLOAD_SCAN_STATUS.PROCESSING) {
      return `/processing-upload`;
    }
    return upload.url;
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
          {this.props.uploads.map(upload => (
            <tr key={upload.id}>
              <td>
                <a href={this.getUploadUrl(upload)} target="_blank" className="usa-link">
                  {upload.filename}
                </a>
              </td>
              <td>{moment(upload.created_at).format('LLL')}</td>
              <td>{bytes(upload.bytes)}</td>
              <td>
                <a href="" onClick={e => this.props.onDelete(e, upload.id)} className="usa-link">
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

export default connect(
  mapStateToProps,
  mapDispatchToProps,
)(UploadsTable);
