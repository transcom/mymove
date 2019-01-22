import React, { Component, Fragment } from 'react';
import { bindActionCreators } from 'redux';
import { connect } from 'react-redux';
import { get, concat, includes, map, reject } from 'lodash';

import { push } from 'react-router-redux';
import { getFormValues, reduxForm, Field } from 'redux-form';

import Alert from 'shared/Alert'; // eslint-disable-line
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import DutyStationSearchBox from 'scenes/ServiceMembers/DutyStationSearchBox';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import Uploader from 'shared/Uploader';
import UploadsTable from 'shared/Uploader/UploadsTable';
import SaveCancelButtons from './SaveCancelButtons';
import { updateOrders, deleteUploads, addUploads } from 'scenes/Orders/ducks';
import { moveIsApproved, isPpm } from 'scenes/Moves/ducks';
import { editBegin, editSuccessful, entitlementChangeBegin, entitlementChanged, checkEntitlement } from './ducks';
import scrollToTop from 'shared/scrollToTop';

import './Review.css';
import profileImage from './images/profile.png';

const editOrdersFormName = 'edit_orders';
const uploaderLabelIdle = 'Drag & drop or <span class="filepond--label-action">click to upload orders</span>';

let EditOrdersForm = props => {
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
  } = props;
  const visibleUploads = reject(existingUploads, upload => {
    return includes(deleteQueue, upload.id);
  });
  return (
    <form onSubmit={handleSubmit}>
      <img src={profileImage} alt="" /> Orders
      <hr />
      <h3 className="sm-heading">Edit Orders:</h3>
      <SwaggerField fieldName="orders_type" swagger={schema} required />
      <SwaggerField fieldName="issue_date" swagger={schema} required />
      <SwaggerField fieldName="report_by_date" swagger={schema} required />
      <SwaggerField fieldName="has_dependents" swagger={schema} component={YesNoBoolean} />
      {get(props, 'formValues.has_dependents', false) && (
        <Fragment>
          <SwaggerField
            fieldName="spouse_has_pro_gear"
            swagger={props.schema}
            component={YesNoBoolean}
            className="wider-label"
          />
        </Fragment>
      )}
      <br />
      <Field name="new_duty_station" component={DutyStationSearchBox} />
      <p>Uploads:</p>
      {Boolean(visibleUploads.length) && <UploadsTable uploads={visibleUploads} onDelete={onDelete} />}
      {Boolean(get(initialValues, 'uploaded_orders')) && (
        <Uploader document={initialValues.uploaded_orders} onChange={onUpload} labelIdle={uploaderLabelIdle} />
      )}
      <SaveCancelButtons valid={valid} submitting={submitting} />
    </form>
  );
};

EditOrdersForm = reduxForm({
  form: editOrdersFormName,
})(EditOrdersForm);

class EditOrders extends Component {
  constructor(props) {
    super(props);

    this.state = {
      newUploads: [],
      deleteQueue: [],
    };

    this.cancelChanges = this.cancelChanges.bind(this);
  }

  cancelChanges = () => {
    const newUploadIds = map(this.state.newUploads, 'id');
    this.props.deleteUploads(newUploadIds).then(() => {
      if (!this.props.hasSubmitError) {
        this.returnToReview();
      } else {
        scrollToTop();
      }
    });
  };

  handleDelete = (e, uploadId) => {
    e.preventDefault();
    this.setState({ deleteQueue: concat(this.state.deleteQueue, uploadId) });
  };

  handleNewUpload = uploads => {
    this.setState({ newUploads: uploads });
  };

  updateOrders = fieldValues => {
    fieldValues.new_duty_station_id = fieldValues.new_duty_station.id;
    fieldValues.spouse_has_pro_gear = (fieldValues.has_dependents && fieldValues.spouse_has_pro_gear) || false;
    let addUploads = this.props.addUploads(this.state.newUploads);
    let deleteUploads = this.props.deleteUploads(this.state.deleteQueue);
    if (
      fieldValues.has_dependents !== this.props.currentOrders.has_dependents ||
      fieldValues.spouse_has_pro_gear !== this.props.spouse_has_pro_gear
    ) {
      this.props.entitlementChanged();
    }
    return Promise.all([addUploads, deleteUploads])
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
              onSubmit={this.updateOrders}
              schema={schema}
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
  const props = {
    currentOrders: state.orders.currentOrders,
    error: get(state, 'orders.error'),
    existingUploads: get(state, `orders.currentOrders.uploaded_orders.uploads`, []),
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
      addUploads,
      deleteUploads,
      editBegin,
      entitlementChangeBegin,
      editSuccessful,
      entitlementChanged,
      checkEntitlement,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(EditOrders);
