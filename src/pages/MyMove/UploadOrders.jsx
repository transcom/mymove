import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import './UploadOrders.css';

import {
  fetchLatestOrders as fetchLatestOrdersAction,
  selectUploadsForActiveOrders,
} from 'shared/Entities/modules/orders';
import {
  createUpload as createUploadAction,
  deleteUpload as deleteUploadAction,
  selectDocument,
} from 'shared/Entities/modules/documents';
import OrdersUploader from 'components/OrdersUploader/index';
import ConnectedUploadsTable from 'shared/Uploader/UploadsTable';
import ConnectedWizardPage from 'shared/WizardPage/index';
import { documentSizeLimitMsg } from 'shared/constants';
import { selectServiceMemberFromLoggedInUser, selectCurrentOrders } from 'store/entities/selectors';
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
    this.deleteFile = this.deleteFile.bind(this);
  }

  componentDidMount() {
    const { serviceMemberId, fetchLatestOrders } = this.props;
    fetchLatestOrders(serviceMemberId);
  }

  onChange(files) {
    this.setState({
      newUploads: files,
    });
  }

  deleteFile(e, uploadId) {
    e.preventDefault();
    const { currentOrders, deleteUpload } = this.props;
    if (currentOrders) {
      deleteUpload(uploadId);
    }
  }

  render() {
    const {
      pages,
      pageKey,
      error,
      currentOrders,
      uploads,
      document,
      additionalParams,
      createUpload,
      deleteUpload,
    } = this.props;
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
            <ConnectedUploadsTable uploads={uploads} onDelete={this.deleteFile} />
          </>
        )}
        {currentOrders && (
          <div className="uploader-box">
            <OrdersUploader
              createUpload={createUpload}
              deleteUpload={deleteUpload}
              document={document}
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
  fetchLatestOrders: PropTypes.func.isRequired,
  createUpload: PropTypes.func.isRequired,
  deleteUpload: PropTypes.func.isRequired,
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
    uploads: selectUploadsForActiveOrders(state),
    document: selectDocument(state, currentOrders?.uploaded_orders),
  };

  return props;
}

const mapDispatchToProps = {
  fetchLatestOrders: fetchLatestOrdersAction,
  createUpload: createUploadAction,
  deleteUpload: deleteUploadAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
