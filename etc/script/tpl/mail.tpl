<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>告警通知</title>
    <style type="text/css">
        .wrapper {
            background-color: #f8f8f8;
            padding: 15px;
            height: 100%;
        }
        .main {
            width: 600px;
            padding: 30px;
            margin: 0 auto;
            background-color: #fff;
            font-size: 12px;
            font-family: verdana,'Microsoft YaHei',Consolas,'Deja Vu Sans Mono','Bitstream Vera Sans Mono';
        }
        header {
            border-radius: 2px 2px 0 0;
        }
        header .title {
            font-size: 16px;
            color: #333333;
            margin: 0;
        }
        header .sub-desc {
            color: #333;
            font-size: 14px;
            margin-top: 6px;
            margin-bottom: 0;
        }
        hr {
            margin: 20px 0;
            height: 0;
            border: none;
            border-top: 1px solid #e5e5e5;
        }
        em {
            font-weight: 600;
        }
        table {
            margin: 20px 0;
            width: 100%;
        }

        table tbody tr{
            font-weight: 200;
            font-size: 12px;
            color: #666;
            height: 32px;
        }

        .succ {
            background-color: green;
            color: white;
        }

        .fail {
            background-color: red;
            color: white;
        }

        table tbody tr th {
            width: 80px;
            text-align: right;
        }
        .text-right {
            text-align: right;
        }
        .body {
            margin-top: 24px;
        }
        .body-text {
            color: #666666;
            -webkit-font-smoothing: antialiased;
        }
        .body-extra {
            -webkit-font-smoothing: antialiased;
        }
        .body-extra.text-right a {
            text-decoration: none;
            color: #333;
        }
        .body-extra.text-right a:hover {
            color: #666;
        }
        .button {
            width: 200px;
            height: 50px;
            margin-top: 20px;
            text-align: center;
            border-radius: 2px;
            background: #2D77EE;
            line-height: 50px;
            font-size: 20px;
            color: #FFFFFF;
            cursor: pointer;
        }
        .button:hover {
            background: rgb(25, 115, 255);
            border-color: rgb(25, 115, 255);
            color: #fff;
        }
        footer {
            margin-top: 10px;
            text-align: right;
        }
        .footer-logo {
            text-align: right;
        }
        .footer-logo-image {
            width: 108px;
            height: 27px;
            margin-right: 10px;
        }
        .copyright {
            margin-top: 10px;
            font-size: 12px;
            text-align: right;
            color: #999;
            -webkit-font-smoothing: antialiased;
        }
    </style>
</head>
<body>
<div class="wrapper">
    <div class="main">
        <header>
            <h3 class="title">{{Sname}}</h3>
            <p class="sub-desc"></p>
        </header>

        <hr>

        <div class="body">
            <table cellspacing="0" cellpadding="0" border="0">
                <tbody>
                % if IsAlert:
                <tr class="fail">
                    <th>级别状态：</th>
                    <td>{{Status}}</td>
                </tr>
                % else:
                <tr class="succ">
                    <th>级别状态：</th>
                    <td>{{Status}}</td>
                </tr>
                % end

                % if IsMachineDep:
                <tr>
                    <th>告警设备：</th>
                    <td>{{Ident}}</td>
                </tr>
                <tr>
                    <th>所属分组：</th>
                    <td>
                            {{Classpath}}<br />
                    </td>
                </tr>
                % end
                <tr>
                    <th>监控指标：</th>
                    <td>{{Metric}}</td>
                </tr>
                <tr>
                    <th>tags：</th>
                    <td>{{Tags}}</td>
                </tr>
                <tr>
                    <th>当前值：</th>
                    <td>{{Value}}</td>
                </tr>
                <tr>
                    <th>报警说明：</th>
                    <td>
                        {{ReadableExpression}}
                    </td>
                </tr>
                <tr>
                    <th>触发时间：</th>
                    <td>
                        {{TriggerTime}}
                    </td>
                </tr>
                <tr>
                    <th>报警详情：</th>
                    <td>{{Elink}}</td>
                </tr>
                <tr>
                    <th>报警策略：</th>
                    <td>{{Slink}}</td>
                </tr>
                </tbody>
            </table>

            <hr>

            <footer>
                <div class="footer-logo">
                </div>
            </footer>
        </div>
    </div>
</div>
</body>
</html>
