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
      showAmendedOrders: false,
    };

    this.onChange = this.onChange.bind(this);
    this.deleteFile = this.deleteFile.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.setShowAmendedOrders = this.setShowAmendedOrders.bind(this);
  }

  componentDidMount() {
    // If we have a logged in user at mount time, do our loading then.
    if (this.props.currentServiceMember) {
      const serviceMemberID = this.props.currentServiceMember.id;
      this.props.showCurrentOrders(serviceMemberID);
    }
  }

  componentDidUpdate(prevProps, prevState) {
    // If we don't have a service member yet, fetch one when loggedInUser loads.
    if (
      !prevProps.currentServiceMember &&
      this.props.currentServiceMember &&
      !this.props.currentOrders
    ) {
      const serviceMemberID = this.props.currentServiceMember.id;
      this.props.showCurrentOrders(serviceMemberID);
    }
  }

  handleSubmit() {
    this.props.addUploads(this.state.newUploads);
  }

  setShowAmendedOrders(show) {
    this.setState({ showAmendedOrders: show });
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
        <div>
          <h1 className="sm-heading">Upload Your Orders</h1>
          <p>
            In order to schedule your move, we need to have a complete copy of
            your orders.
          </p>
          <p>
            You can upload a PDF, or you can take a picture of each page and
            upload the images.
          </p>
        </div>
        {!!uploads.length && (
          <div>
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
        {currentOrders && (
          <div className="uploader-box">
            <Uploader
              document={currentOrders.uploaded_orders}
              onChange={this.onChange}
            />
            <div className="hint">(Each page must be clear and legible)</div>
          </div>
        )}

        {/* TODO: Uncomment when we support upload of amended orders */}
        {/* <div className="amended-orders">
          <p>
            Do you have amended orders? If so, you need to upload those as well.
          </p>
          <YesNoBoolean
            value={showAmendedOrders}
            onChange={this.setShowAmendedOrders}
          />
          {this.state.showAmendedOrders && (
            <div className="uploader-box">
              <h4>Upload amended orders</h4>
              <Uploader document={{}} onChange={no_op} />
              <div className="hint">(Each page must be clear and legible)</div>
            </div>
          )}
        </div> */}
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
    currentServiceMember: get(
      state,
      'loggedInUser.loggedInUser.service_member',
    ),
    currentOrders: state.orders.currentOrders,
    uploads: get(state, 'orders.currentOrders.uploaded_orders.uploads', []),
    user: state.loggedInUser,
    ...state.orders,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
