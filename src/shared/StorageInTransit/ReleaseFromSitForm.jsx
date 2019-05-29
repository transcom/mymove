import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

export class ReleaseFromSitForm extends Component {
  render() {
    const { storageInTransitSchema, minDate, onSubmit } = this.props;

    const minReleaseFromSitDate = new Date(minDate);
    const utcMinDate = new Date(
      minReleaseFromSitDate.getUTCFullYear(),
      minReleaseFromSitDate.getUTCMonth(),
      minReleaseFromSitDate.getUTCDate(),
    );
    const disabledDaysForDayPicker = [{ before: utcMinDate }];

    return (
      <form onSubmit={this.props.handleSubmit(onSubmit)} className="release-from-sit-form">
        <div className="editable-panel-column">
          <SwaggerField
            className="release-from-sit-field"
            fieldName="released_on"
            swagger={storageInTransitSchema}
            title="Released on"
            onChange={this.onChange}
            minDate={minDate}
            disabledDays={disabledDaysForDayPicker}
            required
          />
        </div>
      </form>
    );
  }
}

ReleaseFromSitForm.propTypes = {
  storageInTransitSchema: PropTypes.object.isRequired,
  onSubmit: PropTypes.func.isRequired,
  minDate: PropTypes.instanceOf(Date).isRequired,
};

export const formName = 'release_from_sit_form';
ReleaseFromSitForm = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(ReleaseFromSitForm);

function mapStateToProps(state) {
  return {
    storageInTransitSchema: get(state, 'swaggerPublic.spec.definitions.StorageInTransitReleasePayload', {}), //released_on
  };
}

export default connect(mapStateToProps)(ReleaseFromSitForm);
