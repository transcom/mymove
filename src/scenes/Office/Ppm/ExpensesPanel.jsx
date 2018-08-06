import { get } from 'lodash';
import React, { Component } from 'react';
import { formatCents } from 'shared/formatters';
import { getTabularExpenses } from 'scenes/Office/Ppm/ducks';
import { connect } from 'react-redux';

const dollar = cents => (cents ? '$' + formatCents(cents) : null);

class ExpensesPanel extends Component {
  render() {
    const expenseData = {
      categories: [
        {
          category: 'CONTRACTED_EXPENSE',
          payment_methods: {
            GTCC: 600,
          },
          total: 600,
        },
        {
          category: 'RENTAL_EQUIPMENT',
          payment_methods: {
            MIL_PAY: 500,
          },
          total: 500,
        },
        {
          category: 'TOLLS',
          payment_methods: {
            OTHER_DD: 500,
          },
          total: 500,
        },
      ],
      grand_total: {
        payment_method_totals: {
          GTCC: 600,
          MIL_PAY: 500,
          OTHER_DD: 500,
        },
        total: 1600,
      },
    };
    const { schemaMovingExpenseType } = this.props;

    const tabularData = getTabularExpenses(
      expenseData,
      schemaMovingExpenseType,
    );
    return (
      <div className="calculator-panel expense-panel">
        <div className="calculator-panel-title">Expenses</div>
        <div>
          <table cellSpacing={0}>
            <tbody>
              <tr>
                <th>&nbsp;</th>
                <th className="expense-header payment-method" colSpan={3}>
                  Payment Method
                </th>
                <th colSpan={2}>&nbsp;</th>
              </tr>
              <tr>
                <th
                  className="expense-header"
                  width="40%"
                  style={{ textAlign: 'left' }}
                >
                  Items
                </th>
                <th className="expense-header" width="10%">
                  GTCC
                </th>
                <th className="expense-header" width="15%">
                  Other
                </th>
                <th className="expense-header" width="20%">
                  Total
                </th>
                <th className="expense-header" width="15%">
                  &nbsp;
                </th>
              </tr>
              {tabularData.map(row => {
                return (
                  <tr key={row.type}>
                    <td>{row.type}</td>
                    <td align="right">{dollar(row.GTCC)}</td>
                    <td align="right"> {dollar(row.other)} </td>
                    <td align="right">{dollar(row.total)} </td>
                    <td>&nbsp;</td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </div>
    );
  }
}
function mapStateToProps(state) {
  return {
    schemaMovingExpenseType: get(
      state,
      'swagger.spec.definitions.MovingExpenseType',
      {},
    ),
  };
}
export default connect(mapStateToProps)(ExpensesPanel);
