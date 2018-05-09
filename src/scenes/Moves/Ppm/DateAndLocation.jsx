import { pick, get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Fragment } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import {
  renderField,
  recursivelyAnnotateRequiredFields,
  validateRequiredFields,
  addUiSchemaRequiredFields,
} from 'shared/JsonSchemaForm';
import WizardPage from 'shared/WizardPage';
import YesNoBoolean from 'shared/Inputs/YesNoBoolean';
import { loadPpm, createOrUpdatePpm } from './ducks';

const formName = 'ppp_date_and_location';
const uiSchema = {
  requiredFields: ['planned_move_date', 'pickup_zip', 'destination_zip'],
};
const subsetOfFields = [
  'planned_move_date',
  'pickup_zip',
  'destination_zip',
  'additional_pickup_zip',
  'days_in_storage',
];

export class DateAndLocation extends React.Component {
  state = { showAdditionalPickup: false, showTempStorage: false };
  // static getDerivedStateFromProps(nextProps, prevState) {
  //   const result = {};
  //   if (nextProps.additional_pickup_zip) result.showAdditionalPickup = true;
  //   if (nextProps.days_in_storage > 0) result.showTempStorage = true;
  //   return result;
  // }
  setShowAdditionalPickup = value => {
    this.setState({ showAdditionalPickup: value });
  };
  setShowTempStorage = value => {
    this.setState({ showTempStorage: value });
  };
  componentDidMount() {
    const moveId = this.props.match.params.moveId;
    this.props.loadPpm(moveId);
  }
  handleSubmit = () => {
    const { createOrUpdatePpm, dirty } = this.props;
    const moveId = this.props.match.params.moveId;
    const pendingValues = this.props.formData.values;
    if (dirty) {
      //don't update a ppm unless the size has changed
      createOrUpdatePpm(moveId, pendingValues);
    }
  };
  render() {
    const {
      schema,
      pages,
      pageKey,
      valid,
      dirty,
      hasSubmitSuccess,
      error,
    } = this.props;
    addUiSchemaRequiredFields(schema, uiSchema);
    recursivelyAnnotateRequiredFields(schema);
    const fields = schema.properties || {};
    return (
      <WizardPage
        handleSubmit={this.handleSubmit}
        isAsync={true}
        pageList={pages}
        pageKey={pageKey}
        pageIsValid={valid}
        pageIsDirty={dirty}
        hasSucceeded={hasSubmitSuccess}
        error={error}
      >
        <form>
          <h1 className="sm-heading">PPM Dates & Locations</h1>
          <h3> Move Date </h3>
          {renderField('planned_move_date', fields, '')}
          <h3>Pickup Location</h3>
          {renderField('pickup_zip', fields, '')}
          <p>Do you have stuff at another pickup location?</p>
          <YesNoBoolean
            value={this.state.showAdditionalPickup}
            onChange={this.setShowAdditionalPickup}
          />
          {this.state.showAdditionalPickup && (
            <Fragment>
              {renderField('additional_pickup_zip', fields, '')}
              <p>Making additional stops may decrease your PPM incentive.</p>
            </Fragment>
          )}
          <h3>Destination Location</h3>
          <p>
            Enter the ZIP for your new home if you know it, or for destination
            duty station if you don't
          </p>
          {renderField('destination_zip', fields, '')}
          <p>
            Are you going to put your stuff in temporary storage before moving
            into your new home?
          </p>
          <YesNoBoolean
            value={this.state.showTempStorage}
            onChange={this.setShowTempStorage}
          />
          {this.state.showTempStorage && (
            <Fragment>
              <p>How many days do you plan to put your stuff in storage?</p>
              <p className="explanatory">You can choose up to 90 days.</p>
              {renderField('days_in_storage', fields, '')}
            </Fragment>
          )}
        </form>
      </WizardPage>
    );
  }
}

DateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
  loadPpm: PropTypes.func.isRequired,
  createOrUpdatePpm: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
  hasSubmitSuccess: PropTypes.bool.isRequired,
};

function mapStateToProps(state) {
  const props = {
    schema: {},
    ...state.ppm,
    formData: state.form[formName],
    enableReinitialize: true,
  };
  const defaultPickupZip = get(
    state.loggedInUser,
    'loggedInUser.service_member.residential_address.postal_code',
  );
  props.initialValues = props.currentPpm
    ? pick(props.currentPpm, subsetOfFields)
    : defaultPickupZip
      ? {
          pickup_zip: defaultPickupZip,
        }
      : null;
  if (state.swagger.spec) {
    props.schema =
      state.swagger.spec.definitions.UpdatePersonallyProcuredMovePayload;
  }
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadPpm, createOrUpdatePpm }, dispatch);
}

const DateAndLocationForm = reduxForm({
  form: formName,
  validate: validateRequiredFields,
})(DateAndLocation);
export default connect(mapStateToProps, mapDispatchToProps)(
  DateAndLocationForm,
);

/*
          <p>
            Are you going to put your stuff in temporary storage before moving
            into your new home?
          </p>
          <label>
            <input
              type="radio"
              checked={this.state.showTempStorage}
              onChange={this.setShowTempStorage}
            />Yes
          </label>
          <label>
            <input
              type="radio"
              checked={!this.state.showTempStorage}
              onChange={this.setShowTempStorage}
            />No
          </label>{' '}
          {this.state.showTempStorage && (
            <Fragment>
              <p>How many days do you plan to put your stuff in storage?</p>
              <p className="explanatory">You can choose up to 90 days.</p>
              {renderField('days_in_storage', fields, '')}
            </Fragment>
          )}
*/
