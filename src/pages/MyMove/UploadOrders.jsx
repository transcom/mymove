import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import './UploadOrders.css';

import { selectDocument } from 'shared/Entities/modules/documents';
import OrdersUploader from 'components/OrdersUploader/index';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import ConnectedWizardPage from 'shared/WizardPage/index';
import { documentSizeLimitMsg } from 'shared/constants';
import { getOrdersForServiceMember, createUploadForDocument, deleteUpload } from 'services/internalApi';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectUploadsForCurrentOrders,
} from 'store/entities/selectors';
// eslint-disable-next-line camelcase
import { no_op as noop } from 'shared/utils';
import {
  PageListShape,
  PageKeyShape,
  AdditionalParamsShape,
  OrdersShape,
  UploadsShape,
  DocumentShape,
} from 'types/customerShapes';

const uploaderLabelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload orders</span>';

export class UploadOrders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
    };

    this.onChange = this.onChange.bind(this);
    this.handleUploadFile = this.handleUploadFile.bind(this);
    this.handleDeleteFile = this.handleDeleteFile.bind(this);
  }

  componentDidMount() {
    const { serviceMemberId, updateOrders } = this.props;
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
    });
  }

  handleUploadFile(file) {
    const { document, serviceMemberId, updateOrders } = this.props;
    return createUploadForDocument(file, document?.id).then(() => {
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
      });
    });
  }

  handleDeleteFile(uploadId) {
    const { serviceMemberId, updateOrders } = this.props;

    return deleteUpload(uploadId).then(() => {
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
      });
    });
  }

  onChange(files) {
    this.setState({
      newUploads: files,
    });
  }

  render() {
    const { pages, pageKey, error, currentOrders, uploads, additionalParams } = this.props;
    const { newUploads } = this.state;
    const isValid = Boolean(uploads.length || newUploads.length);
    const isDirty = Boolean(newUploads.length);
    return (
      <ConnectedWizardPage
        additionalParams={additionalParams}
        dirty={isDirty}
        error={error}
        handleSubmit={noop}
        pageIsValid={isValid}
        pageKey={pageKey}
        pageList={pages}
      >
        <div>
          <h1>Upload your orders</h1>
          <p>In order to schedule your move, we need to have a complete copy of your orders.</p>
          <p>You can upload a PDF, or you can take a picture of each page and upload the images.</p>
          <p>{documentSizeLimitMsg}</p>
        </div>
        {Boolean(uploads.length) && (
          <>
            <br />
            <UploadsTable uploads={uploads} onDelete={this.handleDeleteFile} />
          </>
        )}
        {currentOrders && (
          <div className="uploader-box">
            <OrdersUploader
              createUpload={this.handleUploadFile}
              deleteUpload={this.handleDeleteFile}
              onChange={this.onChange}
              options={{ labelIdle: uploaderLabelIdle }}
            />
            <div className="hint">(Each page must be clear and legible.)</div>
          </div>
        )}
      </ConnectedWizardPage>
    );
  }
}

UploadOrders.propTypes = {
  serviceMemberId: PropTypes.string.isRequired,
  updateOrders: PropTypes.func.isRequired,
  pages: PageListShape.isRequired,
  pageKey: PageKeyShape.isRequired,
  currentOrders: OrdersShape,
  error: PropTypes.string,
  uploads: UploadsShape,
  document: DocumentShape,
  additionalParams: AdditionalParamsShape,
};

UploadOrders.defaultProps = {
  currentOrders: null,
  error: null,
  additionalParams: null,
  uploads: [],
  document: null,
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const currentOrders = selectCurrentOrders(state);

  const props = {
    serviceMemberId,
    currentOrders,
    uploads: selectUploadsForCurrentOrders(state),
    document: selectDocument(state, currentOrders?.uploaded_orders),
  };

  return props;
}

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
