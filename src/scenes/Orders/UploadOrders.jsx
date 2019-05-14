import PropTypes from 'prop-types';
import React, { Component, Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';

import { loadServiceMember } from 'scenes/ServiceMembers/ducks';
import { deleteUpload, addUploads } from './ducks';
import Uploader from 'shared/Uploader';
import UploadsTable from 'shared/Uploader/UploadsTable';
import WizardPage from 'shared/WizardPage';

import './UploadOrders.css';

const uploaderLabelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload orders</span>';

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

  handleSubmit() {
    return this.props.addUploads(this.state.newUploads);
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
    const { pages, pageKey, error, currentOrders, uploads } = this.props;
    const isValid = Boolean(uploads.length || this.state.newUploads.length);
    const isDirty = Boolean(this.state.newUploads.length);
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={isValid}
        dirty={isDirty}
        error={error}
      >
        <div>
          <h1 className="sm-heading">Upload Your Orders</h1>
          <p>In order to schedule your move, we need to have a complete copy of your orders.</p>
          <p>You can upload a PDF, or you can take a picture of each page and upload the images.</p>
        </div>
        {Boolean(uploads.length) && (
          <Fragment>
            <br />
            <UploadsTable uploads={uploads} onDelete={this.deleteFile} />
          </Fragment>
        )}
        {currentOrders && (
          <div className="uploader-box">
            <Uploader
              document={currentOrders.uploaded_orders}
              onChange={this.onChange}
              options={{ labelIdle: uploaderLabelIdle }}
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
  deleteUpload: PropTypes.func.isRequired,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadServiceMember, deleteUpload, addUploads }, dispatch);
}
function mapStateToProps(state) {
  const props = {
    uploads: get(state, 'orders.currentOrders.uploaded_orders.uploads', []),
    ...state.orders,
  };
  return props;
}
export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
