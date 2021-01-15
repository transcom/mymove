import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, includes, reject } from 'lodash';

import { push } from 'connected-react-router';
import { getFormValues, reduxForm, Field } from 'redux-form';

import Alert from 'shared/Alert';
import { withContext } from 'shared/AppContext';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import OrdersUploader from 'components/OrdersUploader';
import UploadsTable from 'shared/Uploader/UploadsTable';
import SectionWrapper from 'components/Customer/SectionWrapper';
import SaveCancelButtons from './SaveCancelButtons';

import { updateOrders, fetchLatestOrders } from 'shared/Entities/modules/orders';
import { createUpload, deleteUpload, selectDocument } from 'shared/Entities/modules/documents';
import { editBegin, editSuccessful, entitlementChangeBegin, entitlementChanged, checkEntitlement } from './ducks';
import scrollToTop from 'shared/scrollToTop';
import { documentSizeLimitMsg } from 'shared/constants';
import { createModifiedSchemaForOrdersTypesFlag } from 'shared/featureFlags';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMoveIsApproved,
  selectUploadsForCurrentOrders,
  selectHasCurrentPPM,
} from 'store/entities/selectors';

import './Review.css';
import profileImage from './images/profile.png';
import PropTypes from 'prop-types';

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
  const showAllOrdersTypes = props.context.flags.allOrdersTypes;
  const modifiedSchemaForOrdersTypesFlag = createModifiedSchemaForOrdersTypesFlag(schema);

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
                marginBottom: 16,
                marginTop: 20,
              }}
            >
              Orders
            </h1>
            <SectionWrapper>
              <h2>Edit Orders:</h2>
              <SwaggerField
                fieldName="orders_type"
                swagger={showAllOrdersTypes ? schema : modifiedSchemaForOrdersTypesFlag}
                required
              />
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
            </SectionWrapper>
            <SaveCancelButtons valid={valid} submitting={submitting} />
          </form>
        </div>
      </div>
    </div>
  );
};

EditOrdersForm.propTypes = {
  context: PropTypes.shape({
    flags: PropTypes.shape({
      allOrdersTypes: PropTypes.bool,
    }).isRequired,
  }).isRequired,
};
EditOrdersForm = withContext(
  reduxForm({
    form: editOrdersFormName,
  })(EditOrdersForm),
);

class EditOrders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
      deleteQueue: [],
    };
  }

  handleDelete = (e, uploadId) => {
    e.preventDefault();
    this.props.deleteUpload(uploadId);
  };

  handleNewUpload = (uploads) => {
    this.setState({ newUploads: uploads });
  };

  updateOrders = (fieldValues) => {
    fieldValues.new_duty_station_id = fieldValues.new_duty_station.id;
    fieldValues.spouse_has_pro_gear = (fieldValues.has_dependents && fieldValues.spouse_has_pro_gear) || false;
    if (
      fieldValues.has_dependents !== this.props.currentOrders.has_dependents ||
      fieldValues.spouse_has_pro_gear !== this.props.spouse_has_pro_gear
    ) {
      this.props.entitlementChanged();
    }
    return Promise.all([this.props.updateOrders(fieldValues.id, fieldValues)]).then(() => {
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
  const serviceMember = selectServiceMemberFromLoggedInUser(state);
  const serviceMemberId = serviceMember?.id;
  const currentOrders = selectCurrentOrders(state) || {};
  const uploads = selectUploadsForCurrentOrders(state);

  const props = {
    currentOrders,
    serviceMemberId,
    existingUploads: uploads,
    document: selectDocument(state, currentOrders.uploaded_orders),
    error: get(state, 'orders.error'),
    formValues: getFormValues(editOrdersFormName)(state),
    hasSubmitError: get(state, 'orders.hasSubmitError'),
    moveIsApproved: selectMoveIsApproved(state),
    isPpm: selectHasCurrentPPM(state),
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
