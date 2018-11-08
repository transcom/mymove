import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { get } from 'lodash';
import PropTypes from 'prop-types';
import { reduxForm } from 'redux-form';
import { SwaggerField } from 'shared/JsonSchemaForm/JsonSchemaField';
import { getCookieURL } from './ducks';

const schema = {
  properties: {
    cookie_name: {
      type: 'string',
      title: 'Cookie Name',
    },
  },
};

export class DPSAuthCookie extends Component {
  sendRequest = values => {
    this.props.getCookieURL(values.cookie_name).then(response => {
      console.log(response);
      var cookieURL = get(response, 'payload.cookie_url', '');
      if (cookieURL) {
        window.location = cookieURL;
      }
    });
  };

  render() {
    return (
      <div className="usa-grid">
        <h1 className="sm-heading">Set DPS Auth Cookie</h1>
        <form onSubmit={this.props.handleSubmit(this.sendRequest)}>
          <SwaggerField fieldName="cookie_name" swagger={this.props.schema} />
          <button type="submit">Submit</button>
        </form>
      </div>
    );
  }
}
DPSAuthCookie.propTypes = {
  getCookieURL: PropTypes.func.isRequired,
  schema: PropTypes.object.isRequired,
};

function mapStateToProps(state) {
  return {
    schema,
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ getCookieURL }, dispatch);
}

export default connect(mapStateToProps, mapDispatchToProps)(reduxForm({ form: 'dpsAuthCookie' })(DPSAuthCookie));
