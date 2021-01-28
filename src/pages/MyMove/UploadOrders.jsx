import React, { Component, createRef } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import './UploadOrders.css';

import FileUpload from 'components/FileUpload/FileUpload';
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
import { PageListShape, PageKeyShape, AdditionalParamsShape, OrdersShape, UploadsShape } from 'types/customerShapes';

export class UploadOrders extends Component {
  constructor(props) {
    super(props);

    this.filePondEl = createRef();

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
    const { currentOrders } = this.props;
    const documentId = currentOrders?.uploaded_orders?.id;
    return createUploadForDocument(file, documentId);
  }

  handleUploadComplete() {
    const { serviceMemberId, updateOrders } = this.props;

    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
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

  onChange() {
    this.filePondEl.current?.removeFiles();
    this.handleUploadComplete();
  }

  render() {
    const { pages, pageKey, error, currentOrders, uploads, additionalParams } = this.props;
    const isValid = !!uploads.length;

    return (
      <ConnectedWizardPage
        additionalParams={additionalParams}
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
            <FileUpload
              ref={this.filePondEl}
              createUpload={this.handleUploadFile}
              onChange={this.onChange}
              labelIdle={'Drag & drop or <span class="filepond--label-action">click to upload orders</span>'}
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
  additionalParams: AdditionalParamsShape,
};

UploadOrders.defaultProps = {
  currentOrders: null,
  error: null,
  additionalParams: null,
  uploads: [],
};

function mapStateToProps(state) {
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const currentOrders = selectCurrentOrders(state);

  const props = {
    serviceMemberId,
    currentOrders,
    uploads: selectUploadsForCurrentOrders(state),
  };

  return props;
}

const mapDispatchToProps = {
  updateOrders: updateOrdersAction,
};

export default connect(mapStateToProps, mapDispatchToProps)(UploadOrders);
