import { pick, get } from 'lodash';
import PropTypes from 'prop-types';
import React from 'react';
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

import { loadPpm, createOrUpdatePpm } from './ducks';

const formName = 'ppp_date_and_location';
const uiSchema = {
  requiredFields: ['planned_move_date', 'pickup_zip', 'destination_zip'],
};
const subsetOfFields = ['planned_move_date', 'pickup_zip', 'destination_zip'];
export class DateAndLocation extends React.Component {
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
          <h3>Destination Location</h3>
          <p>
            Enter the ZIP for your new home if you know it, or for destination
            duty station if you don't
          </p>
          {renderField('destination_zip', fields, '')}
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
    schema: get(
      state,
      'swagger.spec.definitions.UpdatePersonallyProcuredMovePayload',
      {},
    ),
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
