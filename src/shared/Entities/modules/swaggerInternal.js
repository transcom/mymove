import { get } from 'lodash';

export function selectDeptIndicatorDisplayKeyValueList(state) {
  return Object.entries(get(state, 'swaggerInternal.spec.definitions.DeptIndicator.x-display-value', {}));
}

export function selectOrdersTypeDisplayKeyValueList(state) {
  return Object.entries(get(state, 'swaggerInternal.spec.definitions.OrdersType.x-display-value', {}));
}

export function selectOrdersTypeDetailDisplayKeyValueList(state) {
  return Object.entries(get(state, 'swaggerInternal.spec.definitions.OrdersTypeDetail.x-display-value', {}));
}
