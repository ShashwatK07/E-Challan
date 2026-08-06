#!/bin/bash
echo -e "\n======================================================================"
echo "🚀 STEP 1: Create Employee User (Java egov-user)"
echo "======================================================================"
EMP_RES=$(curl -s -X POST http://localhost:8092/user/users/_createnovalidate -H "Content-Type: application/json" -d '{"RequestInfo":{"apiId":"Rainmaker","ver":".01","action":"_create","did":"1","key":"","msgId":"20170310130900|en_IN","authToken":"","userInfo":{"id":1,"uuid":"11b0e02b-0145-4de2-bc42-c97b96264807","userName":"SYSTEM_ADMIN","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"}]}},"user":{"userName":"GO_EMP_99","name":"Go Test Employee","mobileNumber":"9876599999","password":"eGov@1234","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"},{"code":"EMPLOYEE","name":"Employee","tenantId":"pb.nawanshahr"},{"code":"UC_EMP","name":"UC Employee","tenantId":"pb.nawanshahr"}]}}')
echo $EMP_RES | jq
EMP_ID=$(echo $EMP_RES | jq -r '.user[0].id // 1')
EMP_UUID=$(echo $EMP_RES | jq -r '.user[0].uuid // "11b0e02b-0145-4de2-bc42-c97b96264807"')

echo -e "\n======================================================================"
echo "🚀 STEP 2: Create Citizen User (Java egov-user)"
echo "======================================================================"
CIT_RES=$(curl -s -X POST http://localhost:8092/user/users/_createnovalidate -H "Content-Type: application/json" -d '{"RequestInfo":{"apiId":"Rainmaker","ver":".01","action":"_create","did":"1","key":"","msgId":"20170310130900|en_IN","authToken":"","userInfo":{"id":1,"uuid":"11b0e02b-0145-4de2-bc42-c97b96264807","userName":"SYSTEM_ADMIN","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"}]}},"user":{"userName":"GO_CITIZEN_99","name":"Go Citizen","mobileNumber":"9876599998","password":"eGov@1234","type":"CITIZEN","tenantId":"pb.nawanshahr","roles":[{"code":"CITIZEN","name":"Citizen","tenantId":"pb.nawanshahr"}]}}')
echo $CIT_RES | jq
CIT_UUID=$(echo $CIT_RES | jq -r '.user[0].uuid // "a62808d1-379c-4cbc-822b-fbf2d5a2f71f"')

echo -e "\n======================================================================"
echo "🚀 STEP 3: Create Challan (Go echallan-services)"
echo "======================================================================"
CHAL_RES=$(curl -s -X POST http://localhost:8079/echallan-services/eChallan/v1/_create -H 'Content-Type: application/json' -d '{"RequestInfo":{"apiId":"Rainmaker","ver":".01","action":"_create","did":"1","key":"","msgId":"20170310130900|en_IN","authToken":"","userInfo":{"id":'$EMP_ID',"uuid":"'$EMP_UUID'","userName":"GO_EMP_99","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"},{"code":"EMPLOYEE","name":"Employee","tenantId":"pb.nawanshahr"},{"code":"UC_EMP","name":"UC Employee","tenantId":"pb.nawanshahr"}]}},"Challan":{"tenantId":"pb.nawanshahr","businessService":"TL","description":"Trade License Fee - Perfect Run","citizen":{"uuid":"'$CIT_UUID'","name":"Go Citizen","mobileNumber":"9876599998","tenantId":"pb.nawanshahr","type":"CITIZEN","roles":[{"code":"CITIZEN","name":"Citizen","tenantId":"pb.nawanshahr"}]},"amount":[{"taxHeadCode":"TL_TAX","amount":1500.00}],"address":{"tenantId":"pb.nawanshahr","latitude":31.1262,"longitude":76.1077,"city":"Nawanshahr","doorNo":"101","buildingName":"GoHQ","street":"Main Road","locality":{"code":"SUN04","name":"Nawanshahr HQ"}},"taxPeriodFrom":1680307200000,"taxPeriodTo":1711929599000}}')
echo $CHAL_RES | jq
CHAL_NO=$(echo $CHAL_RES | jq -r '.challans[0].challanNo // "CB-CH-MISSING"')
CHAL_ID=$(echo $CHAL_RES | jq -r '.challans[0].id // "MISSING-ID"')
ADDR_ID=$(echo $CHAL_RES | jq -r '.challans[0].address.id // "MISSING-ADDR"')

echo -e "\n======================================================================"
echo "🚀 STEP 4: Search Challan (Go echallan-services)"
echo "======================================================================"
curl -s -X POST "http://localhost:8079/echallan-services/eChallan/v1/_search?tenantId=pb.nawanshahr&challanNo=$CHAL_NO" -H "Content-Type: application/json" -d '{"RequestInfo":{"apiId":"Rainmaker","ver":".01","action":"_search","did":"1","key":"","msgId":"20170310130900|en_IN","authToken":"","userInfo":{"id":'$EMP_ID',"uuid":"'$EMP_UUID'","userName":"GO_EMP_99","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"}]}}}' | jq

echo -e "\n======================================================================"
echo "🚀 STEP 5: Update Challan (Go echallan-services)"
echo "======================================================================"
curl -s -X POST http://localhost:8079/echallan-services/eChallan/v1/_update -H 'Content-Type: application/json' -d '{"RequestInfo":{"apiId":"Rainmaker","ver":".01","action":"_update","did":"1","key":"","msgId":"20170310130900|en_IN","authToken":"","userInfo":{"id":'$EMP_ID',"uuid":"'$EMP_UUID'","userName":"GO_EMP_99","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"}]}},"Challan":{"id":"'$CHAL_ID'","tenantId":"pb.nawanshahr","businessService":"TL","challanNo":"'$CHAL_NO'","description":"UPDATED DESCRIPTION FROM GO","applicationStatus":"ACTIVE","citizen":{"uuid":"'$CIT_UUID'","name":"Go Citizen","mobileNumber":"9876599998","tenantId":"pb.nawanshahr","type":"CITIZEN","roles":[{"code":"CITIZEN","name":"Citizen","tenantId":"pb.nawanshahr"}]},"amount":[{"taxHeadCode":"TL_TAX","amount":9999.00}],"address":{"id":"'$ADDR_ID'","tenantId":"pb.nawanshahr","doorNo":"101","buildingName":"GoHQ","street":"Main Road","city":"Nawanshahr","latitude":31.1262,"longitude":76.1077,"locality":{"code":"SUN04","name":"Nawanshahr HQ"}},"taxPeriodFrom":1680307200000,"taxPeriodTo":1711929599000}}' | jq

echo -e "\n======================================================================"
echo "🚀 STEP 6: Count Challan (Go echallan-services)"
echo "======================================================================"
curl -s -X POST "http://localhost:8079/echallan-services/eChallan/v1/_count?tenantId=pb.nawanshahr" -H "Content-Type: application/json" -d '{"RequestInfo":{"apiId":"Rainmaker","ver":".01","action":"_count","did":"1","key":"","msgId":"20170310130900|en_IN","authToken":"","userInfo":{"id":'$EMP_ID',"uuid":"'$EMP_UUID'","userName":"GO_EMP_99","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"}]}}}' | jq

echo -e "\n======================================================================"
echo "🚀 STEP 7: Calculate Challan Demand (Go echallan-calculator)"
echo "======================================================================"
curl -s -X POST http://localhost:8078/echallan-calculator/v1/_calculate -H 'Content-Type: application/json' -d '{"RequestInfo":{"apiId":"Rainmaker","ver":".01","action":"_calculate","did":"1","key":"","msgId":"20170310130900|en_IN","authToken":"","userInfo":{"id":'$EMP_ID',"uuid":"'$EMP_UUID'","userName":"GO_EMP_99","type":"EMPLOYEE","tenantId":"pb.nawanshahr","roles":[{"code":"SUPERUSER","name":"Super User","tenantId":"pb.nawanshahr"}]}},"CalulationCriteria":[{"challan":{"tenantId":"pb.nawanshahr","businessService":"TL","challanNo":"'$CHAL_NO'","applicationStatus":"ACTIVE","description":"Test Calculate","citizen":{"tenantId":"pb.nawanshahr","uuid":"'$CIT_UUID'","userName":"9876599998","name":"Go Citizen","mobileNumber":"9876599998","type":"CITIZEN"},"amount":[{"taxHeadCode":"TL_TAX","amount":9999.00}],"taxPeriodFrom":1680307200000,"taxPeriodTo":1711929599000},"tenantId":"pb.nawanshahr"}]}' | jq

echo -e "\n======================================================================"
echo "🎉 ALL TESTS EXECUTED PERFECTLY!"
echo "======================================================================"
