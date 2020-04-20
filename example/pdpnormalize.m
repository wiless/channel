 x=[1 .5 2 1.5 3]
% x=[0.9036698308249316 0.07680847136348264 0.008695313720607487 0.00367102413012826]
 tt=[0 0.001 0.00210 0.00320 0.015]
% tt= [0 52.625643926984914 173.79154974598615 246.9472076866867] 
ts=min(diff(tt));

stem(tt,x,'r*');
newtt=0:ts:30*ts
grid on;
hold all;
newpdp=zeros(1,length(newtt));
for k=1:length(x)
plot(newtt,x(k)*sinc((newtt-tt(k))/ts))
newpdp=newpdp+x(k)*sinc((newtt-tt(k))/ts);
end
stem(newtt,newpdp,'b--')
xlabel('\tau (s)')
ylabel('Power')
title('PDP and resampling')
legend('original PDP','normalized ts=1e-3')
