// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { forEach } from 'lodash';

import './UploadsTable.css';

import bytes from 'bytes';
import moment from 'moment';

export class UploadsTable extends Component {
  getUploadUrl = upload => {
    let isInfected = false;

    forEach(upload.tags, function(tag) {
      if (tag.key === 'av-status' && tag.value === 'INFECTED') {
        isInfected = true;
      }
    });

    if (isInfected) {
      return `/infected-upload`;
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
                <a href={this.getUploadUrl(upload)} target="_blank">
                  {upload.filename}
                </a>
              </td>
              <td>{moment(upload.created_at).format('LLL')}</td>
              <td>{bytes(upload.bytes)}</td>
              <td>
                <a href="" onClick={e => this.props.onDelete(e, upload.id)}>
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
