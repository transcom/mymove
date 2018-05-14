// import { get, pick } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import editablePanel from './editablePanel';

import { no_op } from 'shared/utils';

// import { updateOrders, loadOrders } from './ducks';
// import { PanelField } from 'shared/EditablePanel';
// import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const OrdersDisplay = props => {
  // const fieldProps = pick(props, ['schema', 'values']);
  return (
    <React.Fragment>
      <div className="editable-panel-column" />
    </React.Fragment>
  );
};

const OrdersEdit = props => {
  // const { schema } = props;
  return (
    <React.Fragment>
      <div className="form-column">
        <label>Orders number</label>
        <input type="text" name="orders-number" />
      </div>
      <div className="form-column">
        <label>Date issued</label>
        <input type="text" name="date-issued" />
      </div>
      <div className="form-column">
        <label>Move type</label>
        <select name="move-type">
          <option value="permanent-change-of-station">
            Permanent Change of Station
          </option>
          <option value="separation">Separation</option>
          <option value="retirement">Retirement</option>
          <option value="local-move">Local Move</option>
          <option value="tdy">Temporary Duty</option>
          <option value="dependent-travel">Dependent Travel</option>
          <option value="bluebark">Bluebark</option>
          <option value="various">Various</option>
        </select>
      </div>
      <div className="form-column">
        <label>Orders type</label>
        <select name="orders-type">
          <option value="shipment-of-hhg-permitted">
            Shipment of HHG Permitted
          </option>
          <option value="pcs-with-tdy-en-route">PCS with TDY En Route</option>
          <option value="shipment-of-hhg-restricted-or-prohibited">
            Shipment of HHG Restricted or Prohibited
          </option>
          <option value="hhg-restricted-area-hhg-prohibited">
            HHG Restricted Area - HHG Prohibited
          </option>
          <option value="course-of-instruction-20-weeks-or-more">
            Course of Instruction 20 Weeks or More
          </option>
          <option value="shipment-of-hhg-prohibited-but-authorized-within-20-weeks">
            Shipment of HHG Prohibited but Authorized within 20 Weeks
          </option>
          <option value="delayed-approval-20-weeks-or-more">
            Delayed Approval 20 Weeks or More
          </option>
        </select>
      </div>
      <div className="form-column">
        <label>Report by</label>
        <input type="date" name="report-by-date" />
      </div>
      <div className="form-column">
        <label>Current duty station</label>
        <input type="text" name="current-duty-station" />
      </div>
      <div className="form-column">
        <label>New duty station</label>
        <input type="text" name="new-duty-station" />
      </div>
      <div className="form-column">
        <b>Entitlements</b>
        <label>Household goods</label>
        <input type="number" name="household-goods-weight" /> lbs
      </div>
      <div className="form-column">
        <label>Pro-gear</label>
        <input type="number" name="pro-gear-weight" /> lbs
      </div>
      <div className="form-column">
        <label>Spouse pro-gear</label>
        <input type="number" name="spouse-pro-gear-weight" /> lbs
      </div>
      <div className="form-column">
        <label>Short-term storage</label>
        <input type="number" name="short-term-storage-days" /> days
      </div>
      <div className="form-column">
        <label>Long-term storage</label>
        <input type="number" name="long-term-storage-days" /> days
      </div>
      <div className="form-column">
        <input
          type="checkbox"
          id="dependents-checkbox"
          name="dependents-authorized"
        />
        <label htmlFor="dependents-checkbox">Dependents authorized</label>
      </div>
    </React.Fragment>
  );
};

const formName = 'office_move_info_orders';

let OrdersPanel = editablePanel(OrdersDisplay, OrdersEdit);
OrdersPanel = reduxForm({ form: formName })(OrdersPanel);

function mapStateToProps(state) {
  return {
    // reduxForm
    formData: state.form[formName],
    initialValues: {},

    // Wrapper
    schema: {},
    hasError: false,
    errorMessage: state.office.error,
    displayValues: {},
    isUpdating: false,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(OrdersPanel);
