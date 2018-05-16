import { get, pick } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';
import editablePanel from './editablePanel';

import { no_op_action } from 'shared/utils';

// import { updateOrders } from './ducks';
import { PanelSwaggerField, PanelField } from 'shared/EditablePanel';

const OrdersDisplay = props => {
  const fieldProps = pick(props, ['schema', 'values']);
  const values = props.values;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelSwaggerField fieldName="orders_number" {...fieldProps} />
        <PanelSwaggerField
          title="Date issued"
          fieldName="issue_date"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Move type"
          fieldName="orders_type"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Orders type"
          fieldName="orders_type_detail"
          {...fieldProps}
        />
        <PanelSwaggerField
          title="Report by"
          fieldName="report_by_date"
          {...fieldProps}
        />
        <PanelField title="Current Duty Station">
          {values.current_duty_station && `${values.current_duty_station.name}`}
        </PanelField>
        <PanelField title="New Duty Station">
          {values.new_duty_station && `${values.new_duty_station.name}`}
        </PanelField>
      </div>
      <div className="editable-panel-column">Entitlements</div>
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

    // Wrapper
    schema: get(state, 'swagger.spec.definitions.Orders'),
    hasError: false,
    errorMessage: state.office.error,
    displayValues: get(state, 'office.officeOrders'),
    isUpdating: false,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      update: no_op_action,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(OrdersPanel);
