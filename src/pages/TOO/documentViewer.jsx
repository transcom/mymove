import React, { Component } from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { GridContainer, Grid } from '@trussworks/react-uswds';
import { selectMoveOrder } from 'shared/Entities/modules/moveOrders';
import {
  selectDeptIndicatorDisplayKeyValueList,
  selectOrdersTypeDetailDisplayKeyValueList,
  selectOrdersTypeDisplayKeyValueList,
} from 'shared/Entities/modules/swaggerInternal';
import { OrdersDetailForm } from 'components/Office/OrdersDetailForm';
import PropTypes from 'prop-types';
import { getMoveOrder } from '../../shared/Entities/modules/moveTaskOrders';
// import classnames from 'classnames';

class DocumentViewer extends Component {
  componentDidMount() {
    // eslint-disable-next-line react/prop-types,react/destructuring-assignment
    const { moveOrderId } = this.props.match.params;
    const { moveOrder } = this.props;
    if (!moveOrder.id) {
      // eslint-disable-next-line react/destructuring-assignment,react/prop-types
      this.props.getMoveOrder(moveOrderId);
    }
  }

  render() {
    const { moveOrder, deptIndicatorOptions, ordersTypeOptions, ordersTypeDetailOptions } = this.props;
    return (
      <GridContainer>
        <Grid row>
          <Grid col gap>
            COL
          </Grid>
          <Grid col>
            {moveOrder.id && (
              <OrdersDetailForm
                initialValues={{
                  currentDutyStation: moveOrder.originDutyStation,
                  newDutyStation: moveOrder.destinationDutyStation,
                  dateIssued: moveOrder.date_issued,
                  reportByDate: moveOrder.report_by_date,
                  departmentIndicator: '',
                  ordersNumber: moveOrder.order_number,
                  ordersType: moveOrder.order_type,
                  ordersTypeDetail: moveOrder.order_type_detail,
                  tac: '',
                  sac: '',
                }}
                ordersTypeOptions={ordersTypeOptions}
                ordersTypeDetailOptions={ordersTypeDetailOptions}
                deptIndicatorOptions={deptIndicatorOptions}
              />
            )}
          </Grid>
        </Grid>
      </GridContainer>
    );
  }
}

DocumentViewer.propTypes = {
  // eslint-disable-next-line react/forbid-prop-types
  moveOrder: PropTypes.object.isRequired,
  deptIndicatorOptions: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.string)).isRequired,
  ordersTypeOptions: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.string)).isRequired,
  ordersTypeDetailOptions: PropTypes.arrayOf(PropTypes.arrayOf(PropTypes.string)).isRequired,
};

const mapStateToProps = (state, ownProps) => {
  const { moveOrderId } = ownProps.match.params;
  const moveOrder = selectMoveOrder(state, moveOrderId);
  const ordersTypeOptions = selectOrdersTypeDisplayKeyValueList(state);
  const ordersTypeDetailOptions = selectOrdersTypeDetailDisplayKeyValueList(state);
  const deptIndicatorOptions = selectDeptIndicatorDisplayKeyValueList(state);

  return { moveOrder, ordersTypeOptions, ordersTypeDetailOptions, deptIndicatorOptions };
};

const mapDispatchToProps = { getMoveOrder };

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(DocumentViewer));
