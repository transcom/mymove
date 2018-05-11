import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import bytes from 'bytes';
import moment from 'moment';

import { loadServiceMember } from 'scenes/ServiceMembers/ducks';
import { showCurrentOrders, deleteUpload, addUploads } from './ducks';
import Uploader from 'shared/Uploader';
import WizardPage from 'shared/WizardPage';

import './UploadOrders.css';

export class UploadOrders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
    };

    this.onChange = this.onChange.bind(this);
    this.deleteFile = this.deleteFile.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidUpdate(prevProps, prevState) {
    // If we don't have a service member yet, fetch one when loggedInUser loads.
    if (
      !prevProps.user.loggedInUser &&
      this.props.user.loggedInUser &&
      !this.props.currentServiceMember
    ) {
      const serviceMemberID = this.props.user.loggedInUser.service_member.id;
      this.props.loadServiceMember(serviceMemberID);
      this.props.showCurrentOrders(serviceMemberID);
    }
  }

  handleSubmit() {
    this.props.addUploads(this.state.newUploads);
  }

  onChange(files) {
    this.setState({
      newUploads: files,
    });
  }

  deleteFile(e, uploadId) {
    e.preventDefault();
    this.props.deleteUpload(uploadId);
  }

  render() {
    const {
      pages,
      pageKey,
      hasSubmitSuccess,
      error,
      currentOrders,
      uploads,
    } = this.props;
    const isValid = Boolean(uploads.length || this.state.newUploads.length);
    const isDirty = Boolean(this.state.newUploads.length);
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        pageIsDirty={isDirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        {!!uploads.length && (
          <div>
            <h1 className="sm-heading">Previous Uploads of Your Orders</h1>
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
                {uploads.map(upload => (
                  <tr key={upload.id}>
                    <td>
                      <a href={upload.url} target="_blank">
                        {upload.filename}
                      </a>
                    </td>
                    <td>{moment(upload.created_at).format('LLL')}</td>
                    <td>{bytes(upload.bytes)}</td>
                    <td>
                      <a href="" onClick={e => this.deleteFile(e, upload.id)}>
                        Delete
                      </a>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        <h1 className="sm-heading">Upload Photos or PDFs of Your Orders</h1>
        {currentOrders && (
          <Uploader
            document={currentOrders.uploaded_orders}
            onChange={this.onChange}
          />
        )}
      </WizardPage>
    );
  }
}

UploadOrders.propTypes = {
  hasSubmitSuccess: PropTypes.bool.isRequired,
  showCurrentOrders: PropTypes.func.isRequired,
  deleteUpload: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    { showCurrentOrders, loadServiceMember, deleteUpload, addUploads },
    dispatch,
  );
}
function mapStateToProps(state) {
  const props = {
    currentOrders: state.orders.currentOrders,
    uploads: get(state, 'orders.currentOrders.uploaded_orders.uploads', []),
    user: state.loggedInUser,
    ...state.orders,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
