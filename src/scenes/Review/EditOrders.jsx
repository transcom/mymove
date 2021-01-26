import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { get } from 'lodash';
import { push } from 'connected-react-router';
import { getFormValues, reduxForm, Field } from 'redux-form';

import Alert from 'shared/Alert';
import { withContext } from 'shared/AppContext';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import FileUpload from 'components/FileUpload/index';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import SectionWrapper from 'components/Customer/SectionWrapper';
import SaveCancelButtons from './SaveCancelButtons';

import { editBegin, editSuccessful, entitlementChangeBegin, entitlementChanged, checkEntitlement } from './ducks';
import scrollToTop from 'shared/scrollToTop';
import { documentSizeLimitMsg } from 'shared/constants';
import { createModifiedSchemaForOrdersTypesFlag } from 'shared/featureFlags';
import { getOrdersForServiceMember, patchOrders, createUploadForDocument, deleteUpload } from 'services/internalApi';
import { updateOrders as updateOrdersAction } from 'store/entities/actions';
import {
  selectServiceMemberFromLoggedInUser,
  selectCurrentOrders,
  selectMoveIsApproved,
  selectUploadsForCurrentOrders,
  selectHasCurrentPPM,
} from 'store/entities/selectors';

import './Review.css';
import profileImage from './images/profile.png';

const editOrdersFormName = 'edit_orders';
const uploaderLabelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload orders</span>';

let EditOrdersForm = (props) => {
  const { onDelete, schema, handleSubmit, submitting, valid, initialValues, existingUploads } = props;
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
              {Boolean(existingUploads.length) && <UploadsTable uploads={existingUploads} onDelete={onDelete} />}
              {Boolean(get(initialValues, 'uploaded_orders')) && (
                <div>
                  <p>{documentSizeLimitMsg}</p>
                  <FileUpload
                    createUpload={props.createUpload}
                    deleteUpload={props.deleteUpload}
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
  handleUploadFile = (file) => {
    const { currentOrders, serviceMemberId, updateOrders } = this.props;
    const documentId = currentOrders?.uploaded_orders?.id;
    return createUploadForDocument(file, documentId).then(() => {
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
      });
    });
  };

  handleDeleteFile = (uploadId) => {
    const { serviceMemberId, updateOrders } = this.props;

    return deleteUpload(uploadId).then(() => {
      getOrdersForServiceMember(serviceMemberId).then((response) => {
        updateOrders(response);
      });
    });
  };

  submitOrders = (fieldValues) => {
    fieldValues.new_duty_station_id = fieldValues.new_duty_station.id;
    fieldValues.spouse_has_pro_gear = (fieldValues.has_dependents && fieldValues.spouse_has_pro_gear) || false;
    if (
      fieldValues.has_dependents !== this.props.currentOrders.has_dependents ||
      fieldValues.spouse_has_pro_gear !== this.props.spouse_has_pro_gear
    ) {
      this.props.entitlementChanged();
    }

    return patchOrders(fieldValues).then((response) => {
      this.props.updateOrders(response);
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

    const { serviceMemberId, updateOrders } = this.props;
    getOrdersForServiceMember(serviceMemberId).then((response) => {
      updateOrders(response);
    });
  }

  render() {
    const { error, schema, currentOrders, formValues, existingUploads, moveIsApproved } = this.props;
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
              onSubmit={this.submitOrders}
              schema={schema}
              createUpload={this.handleUploadFile}
              deleteUpload={this.handleDeleteFile}
              existingUploads={existingUploads}
              onDelete={this.handleDeleteFile}
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
    error: get(state, 'orders.error'),
    formValues: getFormValues(editOrdersFormName)(state),
    hasSubmitError: get(state, 'orders.hasSubmitError'),
    moveIsApproved: selectMoveIsApproved(state),
    isPpm: selectHasCurrentPPM(state),
    schema: get(state, 'swaggerInternal.spec.definitions.CreateUpdateOrders', {}),
  };
  return props;
}

const mapDispatchToProps = {
  push,
  updateOrders: updateOrdersAction,
  editBegin,
  entitlementChangeBegin,
  editSuccessful,
  entitlementChanged,
  checkEntitlement,
};

export default withContext(connect(mapStateToProps, mapDispatchToProps)(EditOrders));
