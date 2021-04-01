import React, { Component, createRef } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { GridContainer, Grid, Alert } from '@trussworks/react-uswds';

import './UploadOrders.css';

import ScrollToTop from 'components/ScrollToTop';
import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import { documentSizeLimitMsg } from 'shared/constants';
import { getOrdersForServiceMember, createUploadForDocument, deleteUpload } from 'services/internalApi';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectUploadsForCurrentOrders,
} from 'store/entities/selectors';
import { OrdersShape, UploadsShape } from 'types/customerShapes';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import WizardNavigation from 'components/Customer/WizardNavigation/WizardNavigation';
import { customerRoutes, generalRoutes } from 'constants/routes';
import formStyles from 'styles/form.module.scss';

export class UploadOrders extends Component {
  constructor(props) {
    super(props);

    this.state = { isLoading: true, serverError: null };

    this.filePondEl = createRef();

    this.onChange = this.onChange.bind(this);
    this.handleUploadFile = this.handleUploadFile.bind(this);
    this.handleDeleteFile = this.handleDeleteFile.bind(this);
  }

  componentDidMount() {
    const { serviceMemberId, updateOrders } = this.props;
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
      this.setState({ isLoading: false });
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
    const { uploads, push } = this.props;
    const isValid = !!uploads.length;

    const handleBack = () => {
      push(customerRoutes.ORDERS_INFO_PATH);
    };
    const handleNext = () => {
      push(generalRoutes.HOME_PATH);
    };

    const { isLoading, serverError } = this.state;
    if (isLoading) return <LoadingPlaceholder />;

    return (
      <GridContainer>
        <ScrollToTop otherDep={serverError} />

        {serverError && (
          <Grid row>
            <Grid col desktop={{ col: 8, offset: 2 }}>
              <Alert type="error" heading="An error occurred">
                {serverError}
              </Alert>
            </Grid>
          </Grid>
        )}

        <Grid row>
          <Grid col desktop={{ col: 8, offset: 2 }}>
            <h1>Upload your orders</h1>
            <p>In order to schedule your move, we need to have a complete copy of your orders.</p>
            <p>You can upload a PDF, or you can take a picture of each page and upload the images.</p>
            <p>{documentSizeLimitMsg}</p>

            {uploads.length > 0 && (
              <>
                <br />
                <UploadsTable uploads={uploads} onDelete={this.handleDeleteFile} />
              </>
            )}

            <div className="uploader-box">
              <FileUpload
                ref={this.filePondEl}
                createUpload={this.handleUploadFile}
                onChange={this.onChange}
                labelIdle={'Drag & drop or <span class="filepond--label-action">click to upload orders</span>'}
              />
              <div className="hint">(Each page must be clear and legible.)</div>
            </div>

            <div className={formStyles.formActions}>
              <WizardNavigation onBackClick={handleBack} disableNext={!isValid} onNextClick={handleNext} />
            </div>
          </Grid>
        </Grid>
      </GridContainer>
    );
  }
}

UploadOrders.propTypes = {
  serviceMemberId: PropTypes.string.isRequired,
  updateOrders: PropTypes.func.isRequired,
  currentOrders: OrdersShape,
  uploads: UploadsShape,
  push: PropTypes.func.isRequired,
};

UploadOrders.defaultProps = {
  currentOrders: null,
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
