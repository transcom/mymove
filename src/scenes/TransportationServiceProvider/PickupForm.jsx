import React from 'react';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { reduxForm } from 'redux-form';

let PickupForm = props => {
  const { schema, onCancel, handleSubmit, submitting, valid } = props;

  return (
    <form onSubmit={handleSubmit}>
      <div className="usa-width-one-whole">
        <h2 className="extras usa-heading">Enter Pack & Pickup</h2>
      </div>
      <div className="editable-panel-column">
        <div className="column-subhead">Dates</div>
        <SwaggerField fieldName="actual_pack_date" swagger={schema} title="Actual packing (first day)" required />
        <SwaggerField fieldName="actual_pickup_date" swagger={schema} title="Actual pickup" required />
      </div>

      <div className="editable-panel-column">
        <div className="column-subhead">Actual weights</div>
        <SwaggerField className="short-field" fieldName="gross_weight" swagger={schema} /> lbs
        <SwaggerField className="short-field" fieldName="tare_weight" swagger={schema} /> lbs
        <SwaggerField title="Net" className="short-field" fieldName="net_weight" swagger={schema} required /> lbs
      </div>

      <button onClick={onCancel}>Cancel</button>
      <button type="submit" disabled={submitting || !valid}>
        Done
      </button>
    </form>
  );
};

PickupForm = reduxForm({ form: 'pickup_shipment' })(PickupForm);

PickupForm.propTypes = {
  schema: PropTypes.object,
  onCancel: PropTypes.func,
  handleSubmit: PropTypes.func,
  submitting: PropTypes.bool,
  valid: PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const { shipment } = state.tsp;
  return {
    ...ownProps,
    initialValues: shipment,
  };
};

export default connect(mapStateToProps)(PickupForm);
