import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, concat, includes, reject } from 'lodash';

import { push } from 'react-router-redux';
import { getFormValues, reduxForm, Field } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { withContext } from 'shared/AppContext';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import OrdersUploader from 'shared/Uploader/OrdersUploader';
import UploadsTable from 'shared/Uploader/UploadsTable';
import SaveCancelButtons from './SaveCancelButtons';
// import { deleteUploads, addUploads } from 'scenes/Orders/ducks';
import {
  updateOrders,
  fetchLatestOrders,
  selectActiveOrders,
  selectUploadsForOrders,
} from 'shared/Entities/modules/orders';
import { createUpload, deleteUpload, selectDocument } from 'shared/Entities/modules/documents';
import { moveIsApproved, isPpm } from 'scenes/Moves/ducks';
import { editBegin, editSuccessful, entitlementChangeBegin, entitlementChanged, checkEntitlement } from './ducks';
import scrollToTop from 'shared/scrollToTop';
import { documentSizeLimitMsg } from 'shared/constants';

import './Review.css';
import profileImage from './images/profile.png';

const editOrdersFormName = 'edit_orders';
const uploaderLabelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload orders</span>';

let EditOrdersForm = (props) => {
  const {
    onDelete,
    onUpload,
    schema,
    handleSubmit,
    submitting,
    valid,
    initialValues,
    existingUploads,
    deleteQueue,
    document,
  } = props;
  const visibleUploads = reject(existingUploads, (upload) => {
    return includes(deleteQueue, upload.id);
  });
  return (
    <div className="grid-container usa-prose">
      <div className="grid-row">
        <div className="grid-col-12">
          <form onSubmit={handleSubmit}>
            <img src={profileImage} alt="" />
            <h1
              style={{
                display: 'inline-block',
                marginLeft: 10,
                marginBottom: 0,
                marginTop: 20,
              }}
            >
              Orders
            </h1>
            <hr />
            <h3 className="sm-heading">Edit Orders:</h3>
            <SwaggerField fieldName="orders_type" swagger={schema} required />
            <SwaggerField fieldName="issue_date" swagger={schema} required />
            <SwaggerField fieldName="report_by_date" swagger={schema} required />
            <SwaggerField fieldName="has_dependents" swagger={schema} component={YesNoBoolean} />
            <br />
            <Field name="new_duty_station" component={DutyStationSearchBox} />
            <p>Uploads:</p>
            {Boolean(visibleUploads.length) && <UploadsTable uploads={visibleUploads} onDelete={onDelete} />}
            {Boolean(get(initialValues, 'uploaded_orders')) && (
              <div>
                <p>{documentSizeLimitMsg}</p>
                <OrdersUploader
                  createUpload={props.createUpload}
                  deleteUpload={props.deleteUpload}
                  document={document}
                  onChange={onUpload}
                  options={{ labelIdle: uploaderLabelIdle }}
                />
              </div>
            )}
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};

EditOrdersForm = reduxForm({
  form: editOrdersFormName,
  enableReinitialize: true,
})(EditOrdersForm);

class EditOrders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
      deleteQueue: [],
    };

    // this.cancelChanges = this.cancelChanges.bind(this);
  }

  // cancelChanges = () => {
  //   const newUploadIds = map(this.state.newUploads, 'id');
  //   this.props.deleteUploads(newUploadIds).then(() => {
  //     if (!this.props.hasSubmitError) {
  //       this.returnToReview();
  //     } else {
  //       scrollToTop();
  //     }
  //   });
  // };

  handleDelete = (e, uploadId) => {
    e.preventDefault();
    this.setState({ deleteQueue: concat(this.state.deleteQueue, uploadId) });
  };

  handleNewUpload = (uploads) => {
    this.setState({ newUploads: uploads });
  };

  updateOrders = (fieldValues) => {
    fieldValues.new_duty_station_id = fieldValues.new_duty_station.id;
    fieldValues.spouse_has_pro_gear = (fieldValues.has_dependents && fieldValues.spouse_has_pro_gear) || false;
    // let addUploads = this.props.addUploads(this.state.newUploads);
    let deleteUploads = this.props.deleteUploads(this.state.deleteQueue);
    if (
      fieldValues.has_dependents !== this.props.currentOrders.has_dependents ||
      fieldValues.spouse_has_pro_gear !== this.props.spouse_has_pro_gear
    ) {
      this.props.entitlementChanged();
    }
    return Promise.all([deleteUploads])
      .then(() => this.props.updateOrders(fieldValues.id, fieldValues))
      .then(() => {
        // This promise resolves regardless of error.
        if (!this.props.hasSubmitError) {
          this.props.editSuccessful();
          this.props.history.goBack();
          if (this.props.isPpm) {
            this.props.checkEntitlement(this.props.match.params.moveId);
          }
        } else {
          scrollToTop();
        }
      });
  };

  componentDidMount() {
    this.props.editBegin();
    this.props.entitlementChangeBegin();
    const { serviceMemberId } = this.props;
    this.props.fetchLatestOrders(serviceMemberId);
  }

  render() {
    const { error, schema, currentOrders, document, formValues, existingUploads, moveIsApproved } = this.props;
    return (
      <div className="usa-grid">
        {error && (
          <div className="usa-width-one-whole error-message">
            <Alert type="error" heading="An error occurred">
              {error.message}
            </Alert>
          </div>
        )}
        {moveIsApproved && (
          <div className="usa-width-one-whole error-message">
            <Alert type="warning" heading="Your move is approved">
              To make a change to your orders, you will need to contact your local PPPO office.
            </Alert>
          </div>
        )}
        {!moveIsApproved && (
          <div className="usa-width-one-whole">
            <EditOrdersForm
              initialValues={currentOrders}
              onSubmit={this.updateOrders}
              document={document}
              schema={schema}
              createUpload={this.props.createUpload}
              deleteUpload={this.props.deleteUpload}
              existingUploads={existingUploads}
              newUploads={this.state.newUploads}
              deleteQueue={this.state.deleteQueue}
              onUpload={this.handleNewUpload}
              onDelete={this.handleDelete}
              formValues={formValues}
            />
          </div>
        )}
      </div>
    );
  }
}

function mapStateToProps(state) {
  const serviceMemberId = get(state, 'serviceMember.currentServiceMember.id');
  const currentOrders = selectActiveOrders(state);
  // const currentOrders = state.orders.currentOrders; // in master
  const uploads = selectUploadsForOrders(state, currentOrders.id);

  const props = {
    currentOrders,
    serviceMemberId: serviceMemberId,
    // existingUploads: get(state, `orders.currentOrders.uploaded_orders.uploads`, []),
    existingUploads: uploads,
    document: selectDocument(state, currentOrders.uploaded_orders),
    error: get(state, 'orders.error'),
    formValues: getFormValues(editOrdersFormName)(state),
    hasSubmitError: get(state, 'orders.hasSubmitError'),
    moveIsApproved: moveIsApproved(state),
    isPpm: isPpm(state),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateUpdateOrders', {}),
  };
  return props;
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      push,
      updateOrders,
      createUpload,
      deleteUpload,
      fetchLatestOrders,
      editBegin,
      entitlementChangeBegin,
      editSuccessful,
      entitlementChanged,
      checkEntitlement,
    },
    dispatch,
  );
}

export default withContext(connect(mapStateToProps, mapDispatchToProps)(EditOrders));
