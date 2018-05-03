import { pick } from 'lodash';
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

import { loadMove } from '../ducks';
import { loadPpm, createOrUpdatePpm } from './ducks';
import { serviceMemberReducer } from '../../ServiceMembers/ducks';

const formName = 'ppp_date_and_location';
const uiSchema = {
  requiredFields: ['planned_move_date', 'pickup_zip', 'destination_zip'],
};
const subsetOfFields = ['planned_move_date', 'pickup_zip', 'destination_zip'];
export class DateAndLocation extends React.Component {
  componentDidMount() {
    //todo: we should make sure this move matches the redux state
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

          {renderField('planned_move_date', fields, '')}
          {renderField('pickup_zip', fields, '')}
          {renderField('destination_zip', fields, '')}
        </form>
      </WizardPage>
    );
  }
}

DateAndLocation.propTypes = {
  schema: PropTypes.object.isRequired,
  loadMove: PropTypes.func.isRequired,
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
  };
  props.initialValues = props.currentPpm
    ? pick(props.currentPpm, subsetOfFields)
    : null; //todo: get pickup zip from service member
  if (state.swagger.spec) {
    props.schema =
      state.swagger.spec.definitions.UpdatePersonallyProcuredMovePayload;
  }
  return props;
}
function mapDispatchToProps(dispatch) {
  return bindActionCreators({ loadMove, loadPpm, createOrUpdatePpm }, dispatch);
}

const DateAndLocationForm = reduxForm({
  form: formName,
  validate: validateRequiredFields,
})(DateAndLocation);
export default connect(mapStateToProps, mapDispatchToProps)(
  DateAndLocationForm,
);
