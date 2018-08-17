// eslint-disable-next-line no-unused-vars
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import PropTypes from 'prop-types';
import { CreateUpload, DeleteUpload } from 'shared/api.js';
import isMobile from 'is-mobile';
import { concat, reject, every, includes } from 'lodash';

import './index.css';

const SurveyDisplay = props => {
  const fieldProps = {
    schema: props.ordersSchema,
    values: props.orders,
  };

  return (
    <React.Fragment>
      <div className="editable-panel-column">
        {props.orders.orders_number ? (
          <PanelField title="Orders Number">
            <Link to={`/moves/${props.move.id}/orders`} target="_blank">
              <SwaggerValue fieldName="orders_number" {...fieldProps} />&nbsp;
              <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
            </Link>
          </PanelField>
        ) : (
          <PanelField title="Orders Number" className="missing">
            missing
            <Link to={`/moves/${props.move.id}/orders`} target="_blank">
              <FontAwesomeIcon className="icon" icon={faExternalLinkAlt} />
            </Link>
          </PanelField>
        )}
        <PanelField
          title="Date issued"
          value={formatDate(props.orders.issue_date)}
        />
        <PanelSwaggerField fieldName="orders_type" {...fieldProps} />
        {props.orders.orders_type_detail ? (
          <PanelSwaggerField fieldName="orders_type_detail" {...fieldProps} />
        ) : (
          <PanelField title="Orders type detail" className="missing">
            missing
          </PanelField>
        )}
        <PanelField
          title="Report by"
          value={formatDate(props.orders.report_by_date)}
        />
        <PanelField title="Current Duty Station">
          {get(props.serviceMember, 'current_station.name', '')}
        </PanelField>
        <PanelField title="New Duty Station">
          {get(props.orders, 'new_duty_station.name', '')}
        </PanelField>
      </div>
      <div className="editable-panel-column">
        {renderEntitlements(props.entitlements, props.orders)}
        {props.orders.has_dependents && (
          <PanelField title="Dependents" value="Authorized" />
        )}
      </div>
    </React.Fragment>
  );
};

const SurveyEdit = props => {
  const schema = props.ordersSchema;
  return (
    <React.Fragment>
      <div className="editable-panel-column">
        <FormSection name="orders">
          <SwaggerField
            fieldName="orders_number"
            swagger={schema}
            className="half-width"
            required
          />
          <SwaggerField
            fieldName="issue_date"
            swagger={schema}
            className="half-width"
          />
          <SwaggerField fieldName="orders_type" swagger={schema} required />
          <SwaggerField
            fieldName="orders_type_detail"
            swagger={schema}
            required
          />
          <SwaggerField fieldName="report_by_date" swagger={schema} />
        </FormSection>

        <FormSection name="serviceMember">
          <div className="usa-input duty-station">
            <Field
              name="current_station"
              component={DutyStationSearchBox}
              props={{ title: 'Current Duty Station' }}
            />
          </div>
        </FormSection>

        <FormSection name="orders">
          <div className="usa-input duty-station">
            <Field
              name="new_duty_station"
              component={DutyStationSearchBox}
              props={{ title: 'New Duty Station' }}
            />
          </div>
        </FormSection>
      </div>
      <div className="editable-panel-column">
        {renderEntitlements(props.entitlements, props.orders)}

        <FormSection name="orders">
          <SwaggerField
            fieldName="has_dependents"
            swagger={schema}
            title="Dependents authorized"
          />
          {get(props, 'formValues.orders.has_dependents', false) && (
            <SwaggerField
              fieldName="spouse_has_pro_gear"
              swagger={schema}
              title="Spouse has pro gear"
            />
          )}
        </FormSection>
      </div>
    </React.Fragment>
  );
};

const formName = 'shipment_pre_move_survey';

let PremoveSurveyPanel = editablePanelify(SurveyDisplay, SurveyEdit);
PremoveSurveyPanel = reduxForm({
  form: formName,
  enableReinitialize: true,
  keepDirtyOnReinitialize: true,
})(PremoveSurveyPanel);

PremoveSurveyPanel.propTypes = {
  shipment: PropTypes.object,
};

function mapStateToProps(state) {
  return {};
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(PremoveSurveyPanel);
