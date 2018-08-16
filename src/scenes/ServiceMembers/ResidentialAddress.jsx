import { get } from 'lodash';
import PropTypes from 'prop-types';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { getFormValues } from 'redux-form';
import { updateServiceMember } from './ducks';
import { reduxifyWizardForm } from 'shared/WizardPage/Form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';

const formName = 'service_member_residential_address';
const ResidentalWizardForm = reduxifyWizardForm(formName);

export class ResidentialAddress extends Component {
  handleSubmit = () => {
    const newAddress = { residential_address: this.props.values };
    return this.props.updateServiceMember(newAddress);
  };

  render() {
    const { pages, pageKey, error, currentServiceMember } = this.props;
    // initialValues has to be null until there are values from the action since only the first values are taken
    const initialValues = get(currentServiceMember, 'residential_address');
    const serviceMemberId = this.props.match.params.serviceMemberId;
    return (
      <ResidentalWizardForm
        handleSubmit={this.handleSubmit}
        className={formName}
        pageList={pages}
        pageKey={pageKey}
        serverError={error}
        initialValues={initialValues}
        additionalParams={{ serviceMemberId }}
      >
        <h1 className="sm-heading">Current Residence Address</h1>
        <SwaggerField
          fieldName="street_address_1"
          swagger={this.props.schema}
          required
        />
        <SwaggerField
          fieldName="street_address_2"
          swagger={this.props.schema}
        />
        <div className="address_inline">
          <SwaggerField
            fieldName="city"
            swagger={this.props.schema}
            className="city_state_zip"
            required
          />
          <SwaggerField
            fieldName="state"
            swagger={this.props.schema}
            className="city_state_zip"
            required
          />
          <SwaggerField
            fieldName="postal_code"
            swagger={this.props.schema}
            className="city_state_zip"
            required
          />
        </div>
      </ResidentalWizardForm>
    );
  }
}
ResidentialAddress.propTypes = {
  schema: PropTypes.object.isRequired,
  updateServiceMember: PropTypes.func.isRequired,
  currentServiceMember: PropTypes.object,
  error: PropTypes.object,
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ updateServiceMember }, dispatch);
}
function mapStateToProps(state) {
  return {
    schema: get(state, 'swagger.spec.definitions.Address', {}),
    values: getFormValues(formName)(state),
    ...state.serviceMember,
  };
}
export default connect(mapStateToProps, mapDispatchToProps)(ResidentialAddress);
