import { get } from 'lodash';
import React from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { reduxForm } from 'redux-form';

import editablePanel from './editablePanel';
import { no_op_action } from 'shared/utils';
import { loadEntitlements } from 'scenes/Office/ducks';

import {
  PanelSwaggerField,
  PanelField,
  SwaggerValue,
} from 'shared/EditablePanel';

import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExternalLinkAlt from '@fortawesome/fontawesome-free-solid/faExternalLinkAlt';

const OrdersDisplay = props => {
  const fieldProps = {
    schema: props.ordersSchema,
    values: props.orders,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <PanelField title="Orders Number">
          <Link to={`/moves/${props.move.id}/orders`} target="_blank">
            <SwaggerValue fieldName="orders_number" {...fieldProps} />&nbsp;
            <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
          </Link>
        </PanelField>
        <PanelSwaggerField
          title="Date issued"
          fieldName="issue_date"
          {...fieldProps}
        />
        <PanelSwaggerField fieldName="orders_type" {...fieldProps} />
        <PanelSwaggerField fieldName="orders_type_detail" {...fieldProps} />
        <PanelSwaggerField
          title="Report by"
          fieldName="report_by_date"
          {...fieldProps}
        />
        <PanelField title="Current Duty Station">
          {get(props.serviceMember, 'current_station.name', '')}
        </PanelField>
        <PanelField title="New Duty Station">
          {get(props.orders, 'new_duty_station.name', '')}
        </PanelField>
      </div>
      <div className="editable-panel-column">
        <span className="editable-panel-column subheader">Entitlements</span>
        <PanelField title="Household Goods">
          {get(props.entitlements, 'total', '').toLocaleString()} lbs
        </PanelField>
        <PanelField title="Pro-gear">
          {get(props.entitlements, 'pro_gear', '').toLocaleString()} lbs
        </PanelField>
        <PanelField title="Spouse pro-gear">
          {get(props.entitlements, 'pro_gear_spouse', '').toLocaleString()} lbs
        </PanelField>
        <PanelField className="Todo" title="Short-term storage">
          90 days
        </PanelField>
        {props.orders.has_dependents && (
          <PanelField title="Dependents" value="Authorized" />
        )}
      </div>
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
    ordersSchema: get(state, 'swagger.spec.definitions.Orders'),

    hasError: false,
    errorMessage: state.office.error,
    orders: get(state, 'office.officeOrders'),
    move: get(state, 'office.officeMove'),
    serviceMember: get(state, 'office.officeServiceMember'),
    entitlements: loadEntitlements(state),
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
