```
/swagger
  api.yaml
 
  /resources
    move-order.yaml  
    payment-request.yaml
    order.yaml
    notifications.yaml
    image-processor.yaml  
  
  /responses  
    responses.yaml
```
paths:
  /move-order-route:
    $ref: '/path/to/payment-request.yaml'
  /payment-invoicing-route:
    $ref: '/path/to/payment-request.yaml'
  /other-payment-invoicing-route:
    $ref: '/path/to/payment-request.yaml'
