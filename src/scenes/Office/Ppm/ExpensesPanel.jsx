import { get } from 'lodash';
import React, { Component } from 'react';
import { bindActionCreators } from 'redux';
import { formatCents } from 'shared/formatters';
import { selectPPMForMove } from 'shared/Entities/modules/ppms';
import { getTabularExpenses, getPpmExpenseSummary } from 'scenes/Office/Ppm/ducks';
import { connect } from 'react-redux';

const dollar = cents => (cents ? '$' + formatCents(cents) : null);

class ExpensesPanel extends Component {
  componentDidMount() {
    if (this.props.ppmId) this.props.getPpmExpenseSummary(this.props.ppmId);
  }
  componentDidUpdate(prevProps) {
    if (this.props.ppmId && this.props.ppmId !== prevProps.ppmId) this.props.getPpmExpenseSummary(this.props.ppmId);
  }
  render() {
    const { schemaMovingExpenseType, expenseData } = this.props;

    const tabularData = getTabularExpenses(expenseData, schemaMovingExpenseType);
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
                <th className="expense-header" width="40%" style={{ textAlign: 'left' }}>
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
function mapStateToProps(state, ownProps) {
  return {
    ppmId: selectPPMForMove(state, ownProps.moveId).id,
    schemaMovingExpenseType: get(state, 'swaggerInternal.spec.definitions.MovingExpenseType', {}),
    expenseData: get(state, 'ppmIncentive.summary'),
  };
}

const mapDispatchToProps = dispatch =>
  bindActionCreators(
    {
      getPpmExpenseSummary,
    },
    dispatch,
  );

export default connect(mapStateToProps, mapDispatchToProps)(ExpensesPanel);
